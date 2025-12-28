package dto

import (
	"encoding/json"

	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
)

type ExecuteEventRequest struct {
	ProjectID uuid.UUID       `json:"projectId"`
	Type      string          `json:"type"`   // e.g. "project.title_changed"
	Status    events.Status   `json:"status"` // optional; often computed
	Input     json.RawMessage `json:"input"`  // event-specific payload
}
