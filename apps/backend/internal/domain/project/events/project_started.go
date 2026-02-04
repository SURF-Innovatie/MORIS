package events

import (
	"context"
	"errors"
	"fmt"
	"time"

	projdomain "github.com/SURF-Innovatie/MORIS/internal/domain/project"
	"github.com/google/uuid"
)

const ProjectStartedType = "project.started"

type ProjectStarted struct {
	Base
	Title           string              `json:"title"`
	Description     string              `json:"description"`
	StartDate       time.Time           `json:"startDate"`
	EndDate         time.Time           `json:"endDate"`
	Members         []projdomain.Member `json:"members_ids"`
	OwningOrgNodeID uuid.UUID           `json:"owning_org_node_id"`
}

func (ProjectStarted) isEvent()     {}
func (ProjectStarted) Type() string { return ProjectStartedType }
func (e ProjectStarted) String() string {
	return fmt.Sprintf("Project started: %s", e.Title)
}

func (e *ProjectStarted) Apply(p *projdomain.Project) {
	p.Title = e.Title
	p.Description = e.Description
	p.StartDate = e.StartDate
	p.EndDate = e.EndDate
	p.OwningOrgNodeID = e.OwningOrgNodeID
	p.Members = e.Members
}

func (e *ProjectStarted) NotificationTemplate() string {
	return "Project proposal '{{event.Title}}' created in '{{org_node.Name}}'"
}

func (e *ProjectStarted) ApprovalRequestTemplate() string {
	return "Project proposal '{{event.Title}}' requires approval."
}

func (e *ProjectStarted) ApprovedTemplate() string {
	return "Project proposal '{{event.Title}}' has been approved."
}

func (e *ProjectStarted) RejectedTemplate() string {
	return "Project proposal '{{event.Title}}' has been rejected."
}

func (e *ProjectStarted) NotificationVariables() map[string]string {
	return map[string]string{
		"event.Title":       e.Title,
		"event.Description": e.Description,
	}
}

func (e *ProjectStarted) RelatedIDs() RelatedIDs {
	return RelatedIDs{OrgNodeID: &e.OwningOrgNodeID}
}

type ProjectStartedInput struct {
	Title           string              `json:"title"`
	Description     string              `json:"description"`
	StartDate       time.Time           `json:"start_date"`
	EndDate         time.Time           `json:"end_date"`
	Members         []projdomain.Member `json:"members_ids"`
	OwningOrgNodeID uuid.UUID           `json:"owning_org_node_id"`
}

func DecideProjectStarted(
	projectID uuid.UUID,
	actor uuid.UUID,
	in ProjectStartedInput,
	status Status,
) (Event, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if in.Title == "" {
		return nil, errors.New("title is required")
	}
	if in.EndDate.Before(in.StartDate) {
		return nil, errors.New("end date before start date")
	}

	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = ProjectStartedMeta.FriendlyName

	return &ProjectStarted{
		Base:            base,
		Title:           in.Title,
		Description:     in.Description,
		StartDate:       in.StartDate,
		EndDate:         in.EndDate,
		Members:         in.Members,
		OwningOrgNodeID: in.OwningOrgNodeID,
	}, nil
}

var ProjectStartedMeta = EventMeta{
	Type:         ProjectStartedType,
	FriendlyName: "Project Proposal",
}

func init() {
	RegisterMeta(ProjectStartedMeta, func() Event {
		return &ProjectStarted{
			Base: Base{FriendlyNameStr: ProjectStartedMeta.FriendlyName},
		}
	})

	RegisterDecider[ProjectStartedInput](ProjectStartedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *projdomain.Project, in ProjectStartedInput, status Status) (Event, error) {
			return DecideProjectStarted(projectID, actor, in, status)
		})

	RegisterInputType(ProjectStartedType, ProjectStartedInput{})
}
