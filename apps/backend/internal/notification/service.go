package notification

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent/notification"
	"github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type Service interface {
	Create(ctx context.Context, p entities.Notification) (*entities.Notification, error)
	Get(ctx context.Context, id uuid.UUID) (*entities.Notification, error)
	Update(ctx context.Context, id uuid.UUID, p entities.Notification) (*entities.Notification, error)
	List(ctx context.Context) ([]entities.Notification, error)
	ListForUser(ctx context.Context, userID uuid.UUID) ([]entities.Notification, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
}

type service struct {
	cli *ent.Client
}

func NewService(cli *ent.Client) Service {
	return &service{cli: cli}
}

func (s *service) Create(ctx context.Context, p entities.Notification) (*entities.Notification, error) {
	row, err := s.cli.Notification.
		Create().
		SetMessage(p.Message).
		SetUser(p.User).
		SetEvent(p.Event).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return (&entities.Notification{}).FromEnt(row, p.User, p.Event), nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*entities.Notification, error) {
	row, err := s.cli.Notification.
		Query().
		Where(notification.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	user, err := row.QueryUser().Only(ctx)
	if err != nil {
		return nil, err
	}

	event, _ := row.QueryEvent().Only(ctx)

	return (&entities.Notification{}).FromEnt(row, user, event), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, p entities.Notification) (*entities.Notification, error) {
	row, err := s.cli.Notification.
		UpdateOneID(id).
		SetMessage(p.Message).
		SetUser(p.User).
		SetEvent(p.Event).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return (&entities.Notification{}).FromEnt(row, p.User, p.Event), nil
}

func (s *service) List(ctx context.Context) ([]entities.Notification, error) {
	rows, err := s.cli.Notification.
		Query().
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]entities.Notification, 0, len(rows))
	for _, r := range rows {
		user, err := r.QueryUser().Only(ctx)
		if err != nil {
			return nil, err
		}

		event, _ := r.QueryEvent().Only(ctx)

		out = append(out, *(&entities.Notification{}).FromEnt(r, user, event))
	}

	return out, nil
}

func (s *service) ListForUser(ctx context.Context, userID uuid.UUID) ([]entities.Notification, error) {
	rows, err := s.cli.Notification.
		Query().
		Where(notification.HasUserWith(user.IDEQ(userID))).
		Order(ent.Desc(notification.FieldSentAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]entities.Notification, 0, len(rows))
	for _, r := range rows {
		u, err := r.QueryUser().Only(ctx)
		if err != nil {
			return nil, err
		}

		// event can be empty
		event, _ := r.QueryEvent().Only(ctx)

		out = append(out, *(&entities.Notification{}).FromEnt(r, u, event))
	}
	return out, nil
}

func (s *service) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	_, err := s.cli.Notification.UpdateOneID(id).SetRead(true).Save(ctx)
	return err
}
