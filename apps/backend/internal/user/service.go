package user

import (
	"context"

	"encoding/json"

	"github.com/SURF-Innovatie/MORIS/ent"
	en "github.com/SURF-Innovatie/MORIS/ent/event"
	entperson "github.com/SURF-Innovatie/MORIS/ent/person"
	"github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	"github.com/SURF-Innovatie/MORIS/internal/person"
	"github.com/google/uuid"
)

type Service interface {
	Get(ctx context.Context, id uuid.UUID) (*entities.User, error)
	Create(ctx context.Context, product entities.User) (*entities.User, error)
	Update(ctx context.Context, id uuid.UUID, product entities.User) (*entities.User, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetAccount(ctx context.Context, id uuid.UUID) (*entities.UserAccount, error)
	GetAccountByEmail(ctx context.Context, email string) (*entities.UserAccount, error)
	GetApprovedEvents(ctx context.Context, userID uuid.UUID) ([]events.Event, error)
	ListAll(ctx context.Context, limit, offset int) ([]*entities.UserAccount, int, error)
	ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error
	SearchPersons(ctx context.Context, query string, observerPersonID *uuid.UUID) ([]entities.Person, error)
}

type service struct {
	cli       *ent.Client
	personSvc person.Service
	es        eventstore.Store
}

func NewService(cli *ent.Client, personSvc person.Service, es eventstore.Store) Service {
	return &service{cli: cli, personSvc: personSvc, es: es}
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	row, err := s.cli.User.
		Query().
		Where(user.IDEQ(id)).
		Only(ctx)

	if err != nil {
		return nil, err
	}

	return (&entities.User{}).FromEnt(row), nil
}

func (s *service) Create(ctx context.Context, user entities.User) (*entities.User, error) {
	// TODO: Validate personID, check password requirements and Hash password before storing it
	row, err := s.cli.User.
		Create().
		SetPersonID(user.PersonID).
		SetPassword(user.Password).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return (&entities.User{}).FromEnt(row), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, user entities.User) (*entities.User, error) {
	row, err := s.cli.User.
		UpdateOneID(id).
		SetPersonID(user.PersonID).
		SetPassword(user.Password).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return (&entities.User{}).FromEnt(row), nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.cli.User.
		DeleteOneID(id).
		Exec(ctx)
}

func (s *service) ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error {
	return s.cli.User.
		UpdateOneID(id).
		SetIsActive(isActive).
		Exec(ctx)
}

func (s *service) GetAccount(ctx context.Context, id uuid.UUID) (*entities.UserAccount, error) {
	userRow, err := s.cli.User.
		Query().
		Where(user.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	personEntity, err := s.personSvc.Get(ctx, userRow.PersonID)
	if err != nil {
		return nil, err
	}

	return &entities.UserAccount{
		User:   *(&entities.User{}).FromEnt(userRow),
		Person: *personEntity,
	}, nil
}

func (s *service) GetAccountByEmail(ctx context.Context, email string) (*entities.UserAccount, error) {
	personEntity, err := s.personSvc.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	userRow, err := s.cli.User.
		Query().
		Where(user.PersonIDEQ(personEntity.ID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return &entities.UserAccount{
		User:   *(&entities.User{}).FromEnt(userRow),
		Person: *personEntity,
	}, nil
}

func (s *service) GetApprovedEvents(ctx context.Context, userID uuid.UUID) ([]events.Event, error) {
	return s.es.LoadUserApprovedEvents(ctx, userID)
}

func (s *service) ListAll(ctx context.Context, limit, offset int) ([]*entities.UserAccount, int, error) {
	total, err := s.cli.User.Query().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	users, err := s.cli.User.Query().
		Limit(limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	accounts := make([]*entities.UserAccount, 0, len(users))
	for _, u := range users {
		acc, err := s.GetAccount(ctx, u.ID)
		if err != nil {
			// Skip users with missing person or other errors for now, or handle appropriately
			continue
		}
		accounts = append(accounts, acc)
	}
	return accounts, total, nil
}

func (s *service) SearchPersons(ctx context.Context, query string, observerPersonID *uuid.UUID) ([]entities.Person, error) {
	// Base query for persons by name or email
	q := s.cli.Person.Query().
		Where(
			entperson.Or(
				entperson.NameContainsFold(query),
				entperson.EmailContainsFold(query),
			),
		)

	// If observer is specified, restrict to persons in shared projects
	if observerPersonID != nil {
		// 1. Find projects where observer is a member
		// Fetch all role assignment events to filter for observer
		allRoleEvents, err := s.cli.Event.
			Query().
			Where(en.TypeEQ(events.ProjectRoleAssignedType)).
			All(ctx)
		if err != nil {
			return nil, err
		}

		projectIDsMap := make(map[uuid.UUID]struct{})
		for _, e := range allRoleEvents {
			b, _ := json.Marshal(e.Data)
			var payload events.ProjectRoleAssigned
			if err := json.Unmarshal(b, &payload); err == nil {
				if payload.PersonID == *observerPersonID {
					projectIDsMap[e.ProjectID] = struct{}{}
				}
			}
		}

		var projectIDs []uuid.UUID
		for id := range projectIDsMap {
			projectIDs = append(projectIDs, id)
		}

		if len(projectIDs) == 0 {
			// No shared projects -> no results (except maybe themselves? optional)
			return []entities.Person{}, nil
		}

		// 2. Find all people in those projects

		memberPersonIDsMap := make(map[uuid.UUID]struct{})
		for _, e := range allRoleEvents {
			if _, ok := projectIDsMap[e.ProjectID]; ok {
				b, _ := json.Marshal(e.Data)
				var payload events.ProjectRoleAssigned
				if err := json.Unmarshal(b, &payload); err == nil {
					memberPersonIDsMap[payload.PersonID] = struct{}{}
				}
			}
		}

		var memberPersonIDs []uuid.UUID
		for id := range memberPersonIDsMap {
			memberPersonIDs = append(memberPersonIDs, id)
		}

		if len(memberPersonIDs) == 0 {
			return []entities.Person{}, nil
		}

		// 3. Filter query by these IDs
		q.Where(entperson.IDIn(memberPersonIDs...))
	}

	// Execute query
	rows, err := q.Limit(20).All(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntities[entities.Person](rows), nil
}
