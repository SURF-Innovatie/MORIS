package entities

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationInfo            NotificationType = "info"
	NotificationApprovalRequest NotificationType = "approval_request"
	NotificationStatusUpdate    NotificationType = "status_update"
)

type Notification struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	EventID *uuid.UUID // optional
	Message string
	Type    NotificationType
	Read    bool
	SentAt  time.Time
}

func (n *Notification) FromEnt(row *ent.Notification) *Notification {
	out := &Notification{
		ID:      row.ID,
		Message: row.Message,
		Type:    NotificationType(row.Type.String()),
		Read:    row.Read,
		SentAt:  row.SentAt,
		UserID:  row.UserID,
	}
	if row.EventID != nil {
		out.EventID = row.EventID
	}
	return out
}
