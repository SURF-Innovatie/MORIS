package user

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/app/person"
	"github.com/SURF-Innovatie/MORIS/internal/domain/identity"
	"github.com/SURF-Innovatie/MORIS/internal/domain/identity/readmodels"
	"github.com/SURF-Innovatie/MORIS/internal/infra/cache"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type Service interface {
	Get(ctx context.Context, id uuid.UUID) (*identity.User, error)
	Create(ctx context.Context, user identity.User) (*identity.User, error)
	Update(ctx context.Context, id uuid.UUID, user identity.User) (*identity.User, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetAccount(ctx context.Context, id uuid.UUID) (*readmodels.UserAccount, error)
	GetAccountByEmail(ctx context.Context, email string) (*readmodels.UserAccount, error)
	ListAll(ctx context.Context, limit, offset int) ([]*readmodels.UserAccount, int, error)
	ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error

	SearchPersons(ctx context.Context, query string, observerPersonID *uuid.UUID) ([]identity.Person, error)
	GetPeopleByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]identity.Person, error)
	GetPeopleByUserIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]identity.Person, error)
}

type service struct {
	users      Repository
	people     person.Service
	membership ProjectMembershipRepository
	cache      cache.UserCache
}

func NewService(
	users Repository,
	people person.Service,
	membership ProjectMembershipRepository,
	cache cache.UserCache,
) Service {
	return &service{users: users, people: people, membership: membership, cache: cache}
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*identity.User, error) {
	if u, err := s.cache.GetUser(ctx, id); err == nil {
		return u, nil
	}
	u, err := s.users.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = s.cache.SetUser(ctx, u)
	return u, nil
}

func (s *service) Create(ctx context.Context, user identity.User) (*identity.User, error) {
	u, err := s.users.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	_ = s.cache.SetUser(ctx, u)
	return u, nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, user identity.User) (*identity.User, error) {
	u, err := s.users.Update(ctx, id, user)
	if err != nil {
		return nil, err
	}
	_ = s.cache.SetUser(ctx, u)
	return u, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.users.Delete(ctx, id); err != nil {
		return err
	}
	return s.cache.DeleteUser(ctx, id)
}

func (s *service) ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error {
	if err := s.users.ToggleActive(ctx, id, isActive); err != nil {
		return err
	}
	return s.cache.DeleteUser(ctx, id)
}

func (s *service) GetAccount(ctx context.Context, id uuid.UUID) (*readmodels.UserAccount, error) {
	u, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	p, err := s.people.Get(ctx, u.PersonID)
	if err != nil {
		return nil, err
	}
	return &readmodels.UserAccount{User: *u, Person: *p}, nil
}

func (s *service) GetAccountByEmail(ctx context.Context, email string) (*readmodels.UserAccount, error) {
	p, err := s.people.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	u, err := s.users.GetByPersonID(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	return &readmodels.UserAccount{User: *u, Person: *p}, nil
}

func (s *service) ListAll(ctx context.Context, limit, offset int) ([]*readmodels.UserAccount, int, error) {
	users, total, err := s.users.ListUsers(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	out := make([]*readmodels.UserAccount, 0, len(users))
	for _, u := range users {
		p, err := s.people.Get(ctx, u.PersonID)
		if err != nil {
			// keep your previous “skip missing person” behavior
			continue
		}
		uu := u // copy
		out = append(out, &readmodels.UserAccount{User: uu, Person: *p})
	}
	return out, total, nil
}

func (s *service) SearchPersons(ctx context.Context, query string, observerPersonID *uuid.UUID) ([]identity.Person, error) {
	// base search
	candidates, err := s.people.Search(ctx, query, 20)
	if err != nil {
		return nil, err
	}
	if observerPersonID == nil {
		return candidates, nil
	}

	projectIDs, err := s.membership.ProjectIDsForPerson(ctx, *observerPersonID)
	if err != nil {
		return nil, err
	}
	if len(projectIDs) == 0 {
		return []identity.Person{}, nil
	}

	memberIDs, err := s.membership.PersonIDsForProjects(ctx, projectIDs)
	if err != nil {
		return nil, err
	}
	if len(memberIDs) == 0 {
		return []identity.Person{}, nil
	}

	allowed := make(map[uuid.UUID]struct{}, len(memberIDs))
	for _, id := range memberIDs {
		allowed[id] = struct{}{}
	}

	// filter candidates in-memory (cheap, limit=20)
	out := lo.Filter(candidates, func(p identity.Person, _ int) bool {
		_, ok := allowed[p.ID]
		return ok
	})
	return out, nil
}

func (s *service) GetPeopleByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]identity.Person, error) {
	persons, err := s.people.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return lo.SliceToMap(persons, func(p identity.Person) (uuid.UUID, identity.Person) {
		return p.ID, p
	}), nil
}

func (s *service) GetPeopleByUserIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]identity.Person, error) {
	userMap := make(map[uuid.UUID]identity.Person)
	for _, uid := range userIDs {
		u, err := s.users.Get(ctx, uid)
		if err != nil {
			continue
		}
		p, err := s.people.Get(ctx, u.PersonID)
		if err != nil {
			continue
		}
		userMap[uid] = *p
	}
	return userMap, nil
}
