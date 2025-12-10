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

// Friendly names mapping
var eventFriendlyNames = map[string]string{
	events.ProjectStartedType:      "Project Proposal",
	events.TitleChangedType:        "Title Change",
	events.DescriptionChangedType:  "Description Change",
	events.StartDateChangedType:    "Start Date Change",
	events.EndDateChangedType:      "End Date Change",
	events.OrganisationChangedType: "Organisation Change",
	events.PersonAddedType:         "Person Addition",
	events.PersonRemovedType:       "Person Removal",
	events.ProductAddedType:        "Product Addition",
	events.ProductRemovedType:      "Product Removal",
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

	eventType := e.Type()
	friendlyName, ok := eventFriendlyNames[eventType]
	if !ok {
		friendlyName = eventType // Fallback to raw type if not in map
	}

	msg := fmt.Sprintf("Your request '%s' has been %s.", friendlyName, status)

	_, err = h.Cli.Notification.
		Create().
		SetMessage(msg).
		SetUser(user).
		SetEventID(e.GetID()).
		SetType(notification.TypeStatusUpdate).
		Save(ctx)

	return err
}
