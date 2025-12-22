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

func (e *ProjectStarted) Apply(project *entities.Project) {
	project.Title = e.Title
	project.Description = e.Description
	project.StartDate = e.StartDate
	project.EndDate = e.EndDate
	project.OwningOrgNodeID = e.OwningOrgNodeID
	project.Members = e.Members
}

func (e *ProjectStarted) NotificationMessage() string {
	return fmt.Sprintf("Project '%s' has been started.", e.Title)
}

func (e *ProjectStarted) RelatedIDs() RelatedIDs {
	return RelatedIDs{OrgNodeID: &e.OwningOrgNodeID}
}

func init() {
	RegisterMeta(EventMeta{
		Type:         ProjectStartedType,
		FriendlyName: "Project Proposal",
	}, func() Event { return &ProjectStarted{} })
}
