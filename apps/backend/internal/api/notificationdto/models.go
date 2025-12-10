package notificationdto

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Response struct {
	ID uuid.UUID `json:"id"`
	// TODO: communicate information about the event
	ProjectID uuid.UUID `json:"projectId"`
	EventID   uuid.UUID `json:"eventId"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	Read      bool      `json:"read"`
	SentAt    time.Time `json:"sentAt"`
}

func FromEntity(n entities.Notification) Response {
	projectId := uuid.Nil
	eventId := uuid.Nil
	if n.Event != nil {
		projectId = n.Event.ProjectID
		eventId = n.Event.ID
	}
	return Response{
		ID:        n.Id,
		Message:   n.Message,
		Type:      n.Type,
		Read:      n.Read,
		ProjectID: projectId,
		EventID:   eventId,
		SentAt:    n.SentAt,
	}
}
