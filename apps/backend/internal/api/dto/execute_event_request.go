package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type ExecuteEventRequest struct {
	ProjectID uuid.UUID       `json:"project_id,omitempty"`
	Type      string          `json:"type"`
	Status    string          `json:"status,omitempty"`
	Input     json.RawMessage `json:"input,omitempty"`
}
