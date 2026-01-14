package user

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/app/person"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type Service interface {
	Get(ctx context.Context, id uuid.UUID) (*entities.User, error)
	Create(ctx context.Context, user entities.User) (*entities.User, error)
	Update(ctx context.Context, id uuid.UUID, user entities.User) (*entities.User, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetAccount(ctx context.Context, id uuid.UUID) (*entities.UserAccount, error)
	GetAccountByEmail(ctx context.Context, email string) (*entities.UserAccount, error)
	GetApprovedEvents(ctx context.Context, userID uuid.UUID) ([]events.Event, error)
	ListAll(ctx context.Context, limit, offset int) ([]*entities.UserAccount, int, error)
	ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error

	SearchPersons(ctx context.Context, query string, observerPersonID *uuid.UUID) ([]entities.Person, error)
}

type service struct {
	users      Repository
	people     person.Repository
	es         EventStore
	membership ProjectMembershipRepository
}

func NewService(
	users Repository,
	people person.Repository,
	es EventStore,
	membership ProjectMembershipRepository,
) Service {
	return &service{users: users, people: people, es: es, membership: membership}
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return s.users.Get(ctx, id)
}

func (s *service) Create(ctx context.Context, user entities.User) (*entities.User, error) {
	return s.users.Create(ctx, user)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, user entities.User) (*entities.User, error) {
	return s.users.Update(ctx, id, user)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.users.Delete(ctx, id)
}

func (s *service) ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error {
	return s.users.ToggleActive(ctx, id, isActive)
}

func (s *service) GetAccount(ctx context.Context, id uuid.UUID) (*entities.UserAccount, error) {
	u, err := s.users.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	p, err := s.people.Get(ctx, u.PersonID)
	if err != nil {
		return nil, err
	}
	return &entities.UserAccount{User: *u, Person: *p}, nil
}

func (s *service) GetAccountByEmail(ctx context.Context, email string) (*entities.UserAccount, error) {
	p, err := s.people.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	u, err := s.users.GetByPersonID(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	return &entities.UserAccount{User: *u, Person: *p}, nil
}

func (s *service) GetApprovedEvents(ctx context.Context, userID uuid.UUID) ([]events.Event, error) {
	return s.es.LoadUserApprovedEvents(ctx, userID)
}

func (s *service) ListAll(ctx context.Context, limit, offset int) ([]*entities.UserAccount, int, error) {
	users, total, err := s.users.ListUsers(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	out := make([]*entities.UserAccount, 0, len(users))
	for _, u := range users {
		p, err := s.people.Get(ctx, u.PersonID)
		if err != nil {
			// keep your previous “skip missing person” behavior
			continue
		}
		uu := u // copy
		out = append(out, &entities.UserAccount{User: uu, Person: *p})
	}
	return out, total, nil
}

func (s *service) SearchPersons(ctx context.Context, query string, observerPersonID *uuid.UUID) ([]entities.Person, error) {
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
		return []entities.Person{}, nil
	}

	memberIDs, err := s.membership.PersonIDsForProjects(ctx, projectIDs)
	if err != nil {
		return nil, err
	}
	if len(memberIDs) == 0 {
		return []entities.Person{}, nil
	}

	allowed := make(map[uuid.UUID]struct{}, len(memberIDs))
	for _, id := range memberIDs {
		allowed[id] = struct{}{}
	}

	// filter candidates in-memory (cheap, limit=20)
	out := lo.Filter(candidates, func(p entities.Person, _ int) bool {
		_, ok := allowed[p.ID]
		return ok
	})
	return out, nil
}
