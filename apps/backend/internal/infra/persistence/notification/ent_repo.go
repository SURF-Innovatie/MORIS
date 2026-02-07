package notification

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	entnotification "github.com/SURF-Innovatie/MORIS/ent/notification"
	entuser "github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/notification"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/google/uuid"
)

type EntRepo struct {
	cli *ent.Client
}

func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

func (r *EntRepo) Create(ctx context.Context, n notification.Notification) (*notification.Notification, error) {
	create := r.cli.Notification.Create().
		SetMessage(n.Message).
		SetUserID(n.UserID).
		SetType(entnotification.Type(n.Type)).
		SetDirection(entnotification.Direction(n.Direction)).
		SetDeliveryStatus(entnotification.DeliveryStatus(n.DeliveryStatus))

	if n.EventID != nil && *n.EventID != uuid.Nil {
		create.SetEventID(*n.EventID)
	}

	// LDN/AS2 fields
	if n.ActivityID != nil {
		create.SetActivityID(*n.ActivityID)
	}
	if n.ActivityType != nil {
		create.SetActivityType(*n.ActivityType)
	}
	if n.Payload != nil {
		create.SetPayload(*n.Payload)
	}
	if n.OriginService != nil {
		create.SetOriginService(*n.OriginService)
	}
	if n.TargetService != nil {
		create.SetTargetService(*n.TargetService)
	}

	row, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntityPtr[notification.Notification](row), nil
}

func (r *EntRepo) Get(ctx context.Context, id uuid.UUID) (*notification.Notification, error) {
	row, err := r.cli.Notification.Query().
		Where(entnotification.IDEQ(id)).
		WithEvent().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	entity := transform.ToEntityPtr[notification.Notification](row)
	if row.Edges.Event != nil {
		meta := events.GetMeta(row.Edges.Event.Type)
		fname := meta.FriendlyName
		entity.EventFriendlyName = &fname
	}
	return entity, nil
}

func (r *EntRepo) Update(ctx context.Context, id uuid.UUID, n notification.Notification) (*notification.Notification, error) {
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
	return transform.ToEntityPtr[notification.Notification](row), nil
}

func (r *EntRepo) List(ctx context.Context) ([]notification.Notification, error) {
	rows, err := r.cli.Notification.Query().
		WithEvent().
		All(ctx)
	if err != nil {
		return nil, err
	}

	dtos := transform.ToEntities[notification.Notification](rows)
	for i, row := range rows {
		if row.Edges.Event != nil {
			meta := events.GetMeta(row.Edges.Event.Type)
			fname := meta.FriendlyName
			dtos[i].EventFriendlyName = &fname
		}
	}
	return dtos, nil
}

func (r *EntRepo) ListForUser(ctx context.Context, userID uuid.UUID) ([]notification.Notification, error) {
	rows, err := r.cli.Notification.Query().
		Where(entnotification.HasUserWith(entuser.IDEQ(userID))).
		Order(ent.Desc(entnotification.FieldSentAt)).
		WithEvent().
		All(ctx)
	if err != nil {
		return nil, err
	}

	dtos := transform.ToEntities[notification.Notification](rows)
	for i, row := range rows {
		if row.Edges.Event != nil {
			meta := events.GetMeta(row.Edges.Event.Type)
			fname := meta.FriendlyName
			dtos[i].EventFriendlyName = &fname
		}
	}
	return dtos, nil
}

func (r *EntRepo) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	_, err := r.cli.Notification.UpdateOneID(id).SetRead(true).Save(ctx)
	return err
}

func (r *EntRepo) MarkAsReadByEventID(ctx context.Context, eventID uuid.UUID) error {
	_, err := r.cli.Notification.Update().
		Where(entnotification.EventIDEQ(eventID)).
		SetRead(true).
		Save(ctx)
	return err
}
