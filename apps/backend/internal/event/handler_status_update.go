package event

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/notification"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	notifservice "github.com/SURF-Innovatie/MORIS/internal/notification"
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
	events.ProjectStartedType:        "Project Proposal",
	events.TitleChangedType:          "Title Change",
	events.DescriptionChangedType:    "Description Change",
	events.StartDateChangedType:      "Start Date Change",
	events.EndDateChangedType:        "End Date Change",
	events.OwningOrgNodeChangedType:  "Owning OwningOrgNode Node Change",
	events.ProjectRoleAssignedType:   "Project Role Assignment",
	events.ProjectRoleUnassignedType: "Project Role Unassignment",
	events.ProductAddedType:          "Product Addition",
	events.ProductRemovedType:        "Product Removal",
}

func (h *StatusUpdateNotificationHandler) Handle(ctx context.Context, e events.Event) error {
	u, err := ResolveUser(ctx, h.Cli, e.CreatedByID())
	if err != nil || u == nil {
		return err
	}

	status := e.GetStatus()
	eventType := e.Type()

	friendlyName, ok := eventFriendlyNames[eventType]
	if !ok {
		friendlyName = eventType
	}

	msg := fmt.Sprintf("Your request '%s' has been %s.", friendlyName, status)

	_, err = h.Cli.Notification.
		Create().
		SetMessage(msg).
		SetUser(u).
		SetEventID(e.GetID()).
		SetType(notification.TypeStatusUpdate).
		Save(ctx)

	return err
}
