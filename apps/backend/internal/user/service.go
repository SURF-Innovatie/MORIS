package user

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/SURF-Innovatie/MORIS/internal/api/userdto"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/person"
	"github.com/google/uuid"
)

type Service interface {
	Get(ctx context.Context, id uuid.UUID) (*userdto.Response, error)
	Create(ctx context.Context, product entities.User) (*userdto.Response, error)
	Update(ctx context.Context, id uuid.UUID, product entities.User) (*userdto.Response, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByEmail(ctx context.Context, email string) (*userdto.Response, error)
}

type service struct {
	cli       *ent.Client
	personSvc person.Service
}

func NewService(cli *ent.Client, personSvc person.Service) Service {
	return &service{cli: cli, personSvc: personSvc}
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*userdto.Response, error) {
	row, err := s.cli.User.
		Query().
		Where(user.IDEQ(id)).
		Only(ctx)

	if err != nil {
		return nil, err
	}

	return s.mapRow(ctx, row)
}

func (s *service) Create(ctx context.Context, user entities.User) (*userdto.Response, error) {
	row, err := s.cli.User.
		Create().
		SetPersonID(user.PersonID).
		SetPassword(user.Password).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return s.mapRow(ctx, row)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, user entities.User) (*userdto.Response, error) {
	row, err := s.cli.User.
		UpdateOneID(id).
		SetPersonID(user.PersonID).
		SetPassword(user.Password).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return s.mapRow(ctx, row)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.cli.User.
		DeleteOneID(id).
		Exec(ctx)
}

func (s *service) GetByEmail(ctx context.Context, email string) (*userdto.Response, error) {
	// TODO: very convoluted way to get user by email, refactor later
	personRow, err := s.personSvc.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	userRow, err := s.cli.User.
		Query().
		Where(user.PersonIDEQ(personRow.Id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapRow(ctx, userRow)
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
