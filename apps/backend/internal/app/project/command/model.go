package command

import (
	"encoding/json"

	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/google/uuid"
)

type AvailableEvent struct {
	Type          string
	FriendlyName  string
	NeedsApproval bool
	Allowed       bool
	InputSchema   map[string]any
}

type ExecuteEventRequest struct {
	ProjectID uuid.UUID
	Type      string
	Status    events.Status
	Input     json.RawMessage
}
