package events

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/notification"
	notifservice "github.com/SURF-Innovatie/MORIS/internal/app/notification"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/SURF-Innovatie/MORIS/internal/infra/adapters/user"
)

type StatusUpdateNotificationHandler struct {
	notifier notifservice.Service
	cli      *ent.Client
}

func NewStatusUpdateHandler(cli *ent.Client) *StatusUpdateNotificationHandler {
	return &StatusUpdateNotificationHandler{cli: cli}
}

func (h *StatusUpdateNotificationHandler) Handle(ctx context.Context, e events.Event) error {
	status := e.GetStatus()
	if status != "approved" && status != "rejected" {
		return nil
	}

	u, err := user.ResolveUser(ctx, h.cli, e.CreatedByID())
	if err != nil || u == nil {
		return err
	}

	// status is already retrieved above
	eventType := e.Type()

	meta := events.GetMeta(eventType)
	friendlyName := meta.FriendlyName
	if friendlyName == "" {
		friendlyName = eventType
	}

	msg := fmt.Sprintf("Your request '%s' has been %s.", friendlyName, status)

	_, err = h.cli.Notification.
		Create().
		SetMessage(msg).
		SetUser(u).
		SetEventID(e.GetID()).
		SetType(notification.TypeStatusUpdate).
		Save(ctx)

	return err
}
