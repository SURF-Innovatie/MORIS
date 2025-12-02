package user

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/SURF-Innovatie/MORIS/internal/api/userdto"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/person"
	"github.com/SURF-Innovatie/MORIS/internal/platform/eventstore"
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

	return &entities.User{
		ID:       row.ID,
		PersonID: row.PersonID,
	}, nil
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

	return &entities.User{
		ID:       row.ID,
		PersonID: row.PersonID,
	}, nil
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

	return &entities.User{
		ID:       row.ID,
		PersonID: row.PersonID,
	}, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.cli.User.
		DeleteOneID(id).
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
		User: entities.User{
			ID:       userRow.ID,
			PersonID: userRow.PersonID,
			// Password is omitted for security reasons
		},
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
		Where(user.PersonIDEQ(personEntity.Id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return &entities.UserAccount{
		User: entities.User{
			ID:       userRow.ID,
			PersonID: userRow.PersonID,
			// Password is omitted for security reasons
		},
		Person: *personEntity,
	}, nil
}

func (s *service) mapRow(ctx context.Context, row *ent.User) (*userdto.Response, error) {
	per, err := s.personSvc.Get(ctx, row.PersonID)
	if err != nil {
		return nil, err
	}

	return &userdto.Response{
		ID:         row.ID,
		PersonID:   row.PersonID,
		Email:      per.Email,
		Name:       per.Name,
		ORCiD:      per.ORCiD,
		GivenName:  per.GivenName,
		FamilyName: per.FamilyName,
	}, nil
}

func (s *service) GetApprovedEvents(ctx context.Context, userID uuid.UUID) ([]events.Event, error) {
	return s.es.LoadUserApprovedEvents(ctx, userID)
}
