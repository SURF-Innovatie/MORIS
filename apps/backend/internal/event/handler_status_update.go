package event

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/notification"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	notifservice "github.com/SURF-Innovatie/MORIS/internal/notification"
	"github.com/google/uuid"
)

type StatusUpdateNotificationHandler struct {
	Notifier notifservice.Service
	Cli      *ent.Client
}

func (h *StatusUpdateNotificationHandler) CanHandle(e events.Event) bool {
	status := e.GetStatus()
	return status == "approved" || status == "rejected"
}

func (h *StatusUpdateNotificationHandler) Handle(ctx context.Context, e events.Event) error {
	status := e.GetStatus()
	creatorID := e.CreatedByID()
	if creatorID == uuid.Nil {
		return nil
	}

	user, err := h.Cli.User.Get(ctx, creatorID)
	if err != nil {
		return nil
	}

	msg := fmt.Sprintf("Your request '%s' has been %s.", e.Type(), status)

	_, err = h.Cli.Notification.
		Create().
		SetMessage(msg).
		SetUser(user).
		SetEventID(e.GetID()).
		SetType(notification.TypeStatusUpdate).
		Save(ctx)

	return err
}
