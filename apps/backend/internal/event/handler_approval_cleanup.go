package event

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	enevent "github.com/SURF-Innovatie/MORIS/ent/event"
	"github.com/SURF-Innovatie/MORIS/ent/notification"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
)

type ApprovalCleanupHandler struct {
	Cli *ent.Client
}

func (h *ApprovalCleanupHandler) CanHandle(e events.Event) bool {
	status := e.GetStatus()
	return status == "approved" || status == "rejected"
}

func (h *ApprovalCleanupHandler) Handle(ctx context.Context, e events.Event) error {
	// Mark all "approval_request" notifications for this event as read
	_, err := h.Cli.Notification.
		Update().
		Where(
			notification.And(
				notification.HasEventWith(enevent.ID(e.GetID())),
				notification.TypeEQ(notification.TypeApprovalRequest),
				notification.Read(false),
			),
		).
		SetRead(true).
		Save(ctx)

	return err
}
