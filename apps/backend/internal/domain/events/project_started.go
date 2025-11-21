package events

import (
	"time"

	"github.com/google/uuid"
)

type ProjectStarted struct {
	Base
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	StartDate      time.Time   `json:"startDate"`
	EndDate        time.Time   `json:"endDate"`
	People         []uuid.UUID `json:"people"`
	OrganisationID uuid.UUID   `json:"organisation"`
}

func (ProjectStarted) isEvent()     {}
func (ProjectStarted) Type() string { return ProjectStartedType }
