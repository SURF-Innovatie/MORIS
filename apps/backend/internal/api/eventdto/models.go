package eventdto

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"projectId"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedBy uuid.UUID `json:"createdBy"`
	At        time.Time `json:"at"`
	Details   string    `json:"details"` // Human readable description
}

type Response struct {
	Events []Event `json:"events"`
}
