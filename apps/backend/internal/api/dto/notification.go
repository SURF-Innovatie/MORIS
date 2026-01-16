package dto

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	EventID   *uuid.UUID `json:"event_id,omitempty"`
	ProjectID *uuid.UUID `json:"project_id,omitempty"`

	Message           string    `json:"message"`
	Type              string    `json:"type"`
	Read              bool      `json:"read"`
	SentAt            time.Time `json:"sent_at"`
	EventFriendlyName string    `json:"event_friendly_name,omitempty"`
}

func (r NotificationResponse) FromEntity(n entities.Notification) NotificationResponse {
	resp := NotificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		EventID:   n.EventID,
		ProjectID: n.ProjectID,
		Message:   n.Message,
		Type:      string(n.Type),
		Read:      n.Read,
		SentAt:    n.SentAt,
	}
	if n.EventFriendlyName != nil {
		resp.EventFriendlyName = *n.EventFriendlyName
	}
	return resp
}
