package events

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/notification"
	notifservice "github.com/SURF-Innovatie/MORIS/internal/app/notification"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events/hydrator"
	"github.com/SURF-Innovatie/MORIS/internal/infra/adapters/user"
)

type StatusUpdateNotificationHandler struct {
	notifier notifservice.Service
	cli      *ent.Client
	hydrator *hydrator.Hydrator
}

func NewStatusUpdateHandler(cli *ent.Client, hydrator *hydrator.Hydrator) *StatusUpdateNotificationHandler {
	return &StatusUpdateNotificationHandler{cli: cli, hydrator: hydrator}
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

	msg := h.buildStatusMessage(ctx, e, status)

	_, err = h.cli.Notification.
		Create().
		SetMessage(msg).
		SetUser(u).
		SetEventID(e.GetID()).
		SetType(notification.TypeStatusUpdate).
		Save(ctx)

	return err
}

func (h *StatusUpdateNotificationHandler) buildStatusMessage(ctx context.Context, e events.Event, status events.Status) string {
	// Check if event implements Notifier
	if n, ok := e.(events.Notifier); ok {
		de := h.hydrator.HydrateOne(ctx, e)

		var template string
		if status == events.StatusApproved {
			template = n.ApprovedTemplate()
		} else if status == events.StatusRejected {
			template = n.RejectedTemplate()
		}

		if template != "" {
			vars := events.AddDetailedEventVariables(n.NotificationVariables(), de)
			msg := events.ResolveTemplate(template, vars)
			if msg != "" {
				return msg
			}
		}
	}

	// Fallback to default message
	eventType := e.Type()
	meta := events.GetMeta(eventType)
	friendlyName := meta.FriendlyName
	if friendlyName == "" {
		friendlyName = eventType
	}

	return fmt.Sprintf("Your request '%s' has been %s.", friendlyName, status)
}
