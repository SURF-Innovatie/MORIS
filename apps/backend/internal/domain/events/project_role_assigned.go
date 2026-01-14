package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

const ProjectRoleAssignedType = "project.project_role_assigned"

type ProjectRoleAssigned struct {
	Base
	PersonID      uuid.UUID `json:"person_id"`
	ProjectRoleID uuid.UUID `json:"project_role_id"`
}

func (ProjectRoleAssigned) isEvent()     {}
func (ProjectRoleAssigned) Type() string { return ProjectRoleAssignedType }
func (e ProjectRoleAssigned) String() string {
	return fmt.Sprintf("Role assigned: %s to %s", e.ProjectRoleID, e.PersonID)
}

func (e *ProjectRoleAssigned) Apply(project *entities.Project) {
	for _, m := range project.Members {
		if m.PersonID == e.PersonID && m.ProjectRoleID == e.ProjectRoleID {
			return // already assigned
		}
	}
	project.Members = append(project.Members, entities.ProjectMember{
		PersonID:      e.PersonID,
		ProjectRoleID: e.ProjectRoleID,
	})
}

func (e *ProjectRoleAssigned) RelatedIDs() RelatedIDs {
	return RelatedIDs{PersonID: &e.PersonID, ProjectRoleID: &e.ProjectRoleID}
}

func (e *ProjectRoleAssigned) ApprovalMessage(projectTitle string) string {
	return fmt.Sprintf("Approval requested: role assigned in project '%s'.", projectTitle)
}

type ProjectRoleAssignedInput struct {
	PersonID      uuid.UUID `json:"person_id"`
	ProjectRoleID uuid.UUID `json:"project_role_id"`
}

func DecideProjectRoleAssigned(
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *entities.Project,
	in ProjectRoleAssignedInput,
	status Status,
) (Event, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if in.PersonID == uuid.Nil {
		return nil, errors.New("person id is required")
	}
	if cur == nil {
		return nil, errors.New("current project is required")
	}

	for _, m := range cur.Members {
		if m.PersonID == in.PersonID && m.ProjectRoleID == in.ProjectRoleID {
			return nil, nil
		}
	}

	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = ProjectRoleAssignedMeta.FriendlyName

	return &ProjectRoleAssigned{
		Base:          base,
		PersonID:      in.PersonID,
		ProjectRoleID: in.ProjectRoleID,
	}, nil
}

var ProjectRoleAssignedMeta = EventMeta{
	Type:         ProjectRoleAssignedType,
	FriendlyName: "Project Role Assignment",
}

func init() {
	RegisterMeta(ProjectRoleAssignedMeta, func() Event {
		return &ProjectRoleAssigned{
			Base: Base{FriendlyNameStr: ProjectRoleAssignedMeta.FriendlyName},
		}
	})

	RegisterDecider[ProjectRoleAssignedInput](ProjectRoleAssignedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *entities.Project, in ProjectRoleAssignedInput, status Status) (Event, error) {
			return DecideProjectRoleAssigned(projectID, actor, cur, in, status)
		})

	RegisterInputType(ProjectRoleAssignedType, ProjectRoleAssignedInput{})
}
