package organisation

import (
	"context"

	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/ent"
	or "github.com/SURF-Innovatie/MORIS/ent/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type Service interface {
	Create(ctx context.Context, o entities.Organisation) (*entities.Organisation, error)
	Get(ctx context.Context, id uuid.UUID) (*entities.Organisation, error)
	Update(ctx context.Context, id uuid.UUID, o entities.Organisation) (*entities.Organisation, error)
	List(ctx context.Context) ([]entities.Organisation, error)
}

type service struct {
	cli *ent.Client
}

func NewService(cli *ent.Client) Service {
	return &service{cli: cli}
}

func (s *service) Create(ctx context.Context, o entities.Organisation) (*entities.Organisation, error) {
	row, err := s.cli.Organisation.
		Create().
		SetName(o.Name).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapRow(row), nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*entities.Organisation, error) {
	row, err := s.cli.Organisation.
		Query().
		Where(or.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return mapRow(row), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, p entities.Organisation) (*entities.Organisation, error) {
	row, err := s.cli.Organisation.
		UpdateOneID(id).
		SetName(p.Name).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapRow(row), nil
}

func (s *service) List(ctx context.Context) ([]entities.Organisation, error) {
	rows, err := s.cli.Organisation.
		Query().
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]entities.Organisation, 0, len(rows))
	for _, r := range rows {
		out = append(out, *mapRow(r))
	}
	return out, nil
}

func mapRow(r *ent.Organisation) *entities.Organisation {
	return &entities.Organisation{
		Id:   r.ID,
		Name: r.Name,
	}
}
