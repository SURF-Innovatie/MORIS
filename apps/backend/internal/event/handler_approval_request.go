package event

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/event"
	"github.com/SURF-Innovatie/MORIS/ent/notification"
	"github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	notifservice "github.com/SURF-Innovatie/MORIS/internal/notification"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ApprovalRequestNotificationHandler struct {
	Notifier notifservice.Service
	Cli      *ent.Client
}

func (h *ApprovalRequestNotificationHandler) CanHandle(e events.Event) bool {
	_, ok := e.(events.PersonAdded)
	return ok
}

func (h *ApprovalRequestNotificationHandler) Handle(ctx context.Context, e events.Event) error {
	if _, ok := e.(events.PersonAdded); !ok {
		return nil
	}

	if e.GetStatus() != "pending" {
		return nil
	}

	projectID := e.AggregateID()
	logrus.Infof("NotifyApprovers: Processing PersonAdded event for project %s", projectID)

	startedEvent, err := h.Cli.Event.
		Query().
		Where(
			event.ProjectIDEQ(projectID),
			event.TypeEQ(events.ProjectStartedType),
		).
		First(ctx)
	if err != nil {
		logrus.Errorf("NotifyApprovers: Failed to find ProjectStarted event for project %s: %v", projectID, err)
		return nil
	}

	payload, err := startedEvent.QueryProjectStarted().Only(ctx)
	if err != nil {
		logrus.Errorf("NotifyApprovers: Failed to get ProjectStarted payload: %v", err)
		return nil
	}

	adminID := payload.ProjectAdmin
	if adminID == uuid.Nil {
		logrus.Warnf("NotifyApprovers: Admin ID is nil for project %s", projectID)
		return nil
	}

	adminUser, err := h.Cli.User.Query().
		Where(user.PersonIDEQ(adminID)).
		Only(ctx)
	if err != nil {
		logrus.Errorf("NotifyApprovers: Failed to find user for admin person %s: %v", adminID, err)
		return nil
	}

	msg := fmt.Sprintf("Approval requested: Person added to project '%s'", payload.Title)
	logrus.Infof("NotifyApprovers: Sending notification to admin %s: %s", adminID, msg)

	_, err = h.Cli.Notification.
		Create().
		SetMessage(msg).
		SetUser(adminUser).
		SetEventID(e.GetID()).
		SetType(notification.TypeApprovalRequest).
		Save(ctx)
	return err
}
