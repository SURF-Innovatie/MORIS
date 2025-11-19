package events

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type ProjectStarted struct {
	Base
	Title        string                `json:"title"`
	Description  string                `json:"description"`
	StartDate    time.Time             `json:"startDate"`
	EndDate      time.Time             `json:"endDate"`
	People       []uuid.UUID           `json:"people"`
	Organisation entities.Organisation `json:"organisation"`
}

func (ProjectStarted) isEvent()     {}
func (ProjectStarted) Type() string { return ProjectStartedType }
