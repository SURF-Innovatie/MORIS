package notification

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	entnotification "github.com/SURF-Innovatie/MORIS/ent/notification"
	entuser "github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type EntRepo struct {
	cli *ent.Client
}

func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

func (r *EntRepo) Create(ctx context.Context, n entities.Notification) (*entities.Notification, error) {
	create := r.cli.Notification.Create().
		SetMessage(n.Message).
		SetUserID(n.UserID)

	if n.EventID != nil && *n.EventID != uuid.Nil {
		create.SetEventID(*n.EventID)
	}

	row, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}

	return (&entities.Notification{}).FromEnt(row), nil
}

func (r *EntRepo) Get(ctx context.Context, id uuid.UUID) (*entities.Notification, error) {
	row, err := r.cli.Notification.Query().
		Where(entnotification.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return (&entities.Notification{}).FromEnt(row), nil
}

func (r *EntRepo) Update(ctx context.Context, id uuid.UUID, n entities.Notification) (*entities.Notification, error) {
	upd := r.cli.Notification.UpdateOneID(id).
		SetMessage(n.Message).
		SetUserID(n.UserID)

	if n.EventID != nil && *n.EventID != uuid.Nil {
		upd.SetEventID(*n.EventID)
	} else {
		upd.ClearEvent()
	}

	row, err := upd.Save(ctx)
	if err != nil {
		return nil, err
	}
	return (&entities.Notification{}).FromEnt(row), nil
}

func (r *EntRepo) List(ctx context.Context) ([]entities.Notification, error) {
	rows, err := r.cli.Notification.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]entities.Notification, 0, len(rows))
	for _, row := range rows {
		out = append(out, *(&entities.Notification{}).FromEnt(row))
	}
	return out, nil
}

func (r *EntRepo) ListForUser(ctx context.Context, userID uuid.UUID) ([]entities.Notification, error) {
	rows, err := r.cli.Notification.Query().
		Where(entnotification.HasUserWith(entuser.IDEQ(userID))).
		Order(ent.Desc(entnotification.FieldSentAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]entities.Notification, 0, len(rows))
	for _, row := range rows {
		out = append(out, *(&entities.Notification{}).FromEnt(row))
	}
	return out, nil
}

func (r *EntRepo) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	_, err := r.cli.Notification.UpdateOneID(id).SetRead(true).Save(ctx)
	return err
}
