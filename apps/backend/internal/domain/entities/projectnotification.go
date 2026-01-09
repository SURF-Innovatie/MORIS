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
	ID        uuid.UUID
	UserID    uuid.UUID
	EventID   *uuid.UUID // optional
	ProjectID *uuid.UUID // optional, derived from event
	Message   string
	Type      NotificationType
	Read      bool
	SentAt    time.Time
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
	if row.Edges.Event != nil {
		out.ProjectID = &row.Edges.Event.ProjectID
	}
	return out
}
