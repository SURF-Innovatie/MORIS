package events

import (
	"fmt"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type ProjectStarted struct {
	Base
	Title           string                   `json:"title"`
	Description     string                   `json:"description"`
	StartDate       time.Time                `json:"startDate"`
	EndDate         time.Time                `json:"endDate"`
	Members         []entities.ProjectMember `json:"members_ids"`
	OwningOrgNodeID uuid.UUID                `json:"owning_org_node_id"`
}

func (ProjectStarted) isEvent()     {}
func (ProjectStarted) Type() string { return ProjectStartedType }
func (e ProjectStarted) String() string {
	return fmt.Sprintf("Project started: %s", e.Title)
}
