package person

import (
	"context"

	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/ent"
	pe "github.com/SURF-Innovatie/MORIS/ent/person"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type Service interface {
	Create(ctx context.Context, p entities.Person) (*entities.Person, error)
	Get(ctx context.Context, id uuid.UUID) (*entities.Person, error)
	Update(ctx context.Context, id uuid.UUID, p entities.Person) (*entities.Person, error)
	List(ctx context.Context) ([]*entities.Person, error)
	GetByEmail(ctx context.Context, email string) (*entities.Person, error)
}

type service struct {
	cli *ent.Client
}

func NewService(cli *ent.Client) Service {
	return &service{cli: cli}
}

func (s *service) Create(ctx context.Context, p entities.Person) (*entities.Person, error) {
	row, err := s.cli.Person.
		Create().
		SetName(p.Name).
		SetNillableGivenName(p.GivenName).
		SetNillableFamilyName(p.FamilyName).
		SetEmail(p.Email).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapRow(row), nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*entities.Person, error) {
	row, err := s.cli.Person.
		Query().
		Where(pe.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return mapRow(row), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, p entities.Person) (*entities.Person, error) {
	row, err := s.cli.Person.
		UpdateOneID(id).
		SetName(p.Name).
		SetNillableGivenName(p.GivenName).
		SetNillableFamilyName(p.FamilyName).
		SetEmail(p.Email).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapRow(row), nil
}

func (s *service) List(ctx context.Context) ([]*entities.Person, error) {
	rows, err := s.cli.Person.
		Query().
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*entities.Person, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapRow(r))
	}
	return out, nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (*entities.Person, error) {
	row, err := s.cli.Person.
		Query().
		Where(pe.EmailEQ(email)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return mapRow(row), nil
}

func mapRow(r *ent.Person) *entities.Person {
	return &entities.Person{
		Id:         r.ID,
		Name:       r.Name,
		ORCiD:      &r.OrcidID,
		GivenName:  r.GivenName,
		FamilyName: r.FamilyName,
		Email:      r.Email,
	}
}
