package notificationdto

import (
	"time"

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
