package notification

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent/projectnotification"
	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type Service interface {
	Create(ctx context.Context, p entities.ProjectNotification) (*entities.ProjectNotification, error)
	Get(ctx context.Context, id uuid.UUID) (*entities.ProjectNotification, error)
	Update(ctx context.Context, id uuid.UUID, p entities.ProjectNotification) (*entities.ProjectNotification, error)
	List(ctx context.Context) ([]entities.ProjectNotification, error)
}

type service struct {
	cli *ent.Client
}

func NewService(cli *ent.Client) Service {
	return &service{cli: cli}
}

func (s *service) Create(ctx context.Context, p entities.ProjectNotification) (*entities.ProjectNotification, error) {
	row, err := s.cli.ProjectNotification.
		Create().
		SetMessage(p.Message).
		SetUser(p.User).
		SetProjectID(p.ProjectId).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return mapRow(row, p.User), nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*entities.ProjectNotification, error) {
	row, err := s.cli.ProjectNotification.
		Query().
		Where(projectnotification.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	user, err := row.QueryUser().Only(ctx)
	if err != nil {
		return nil, err
	}
	return mapRow(row, user), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, p entities.ProjectNotification) (*entities.ProjectNotification, error) {
	row, err := s.cli.ProjectNotification.
		UpdateOneID(id).
		SetMessage(p.Message).
		SetUser(p.User).
		SetProjectID(p.ProjectId).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapRow(row, p.User), nil
}

func (s *service) List(ctx context.Context) ([]entities.ProjectNotification, error) {
	rows, err := s.cli.ProjectNotification.
		Query().
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]entities.ProjectNotification, 0, len(rows))
	for _, r := range rows {
		user, err := r.QueryUser().Only(ctx)
		if err != nil {
			return nil, err
		}
		out = append(out, *mapRow(r, user))
	}
	return out, nil
}

func mapRow(r *ent.ProjectNotification, u *ent.User) *entities.ProjectNotification {
	return &entities.ProjectNotification{
		Id:        r.ID,
		ProjectId: r.ProjectID,
		Message:   r.Message,
		User:      u,
	}
}
