package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
	"github.com/google/uuid"
)

const ProjectRoleUnassignedType = "project.role_unassigned"

type ProjectRoleUnassigned struct {
	Base
	PersonID      uuid.UUID `json:"person_id"`
	ProjectRoleID uuid.UUID `json:"project_role_id"`
}

func (ProjectRoleUnassigned) isEvent()     {}
func (ProjectRoleUnassigned) Type() string { return ProjectRoleUnassignedType }
func (e ProjectRoleUnassigned) String() string {
	return fmt.Sprintf("Role unassigned: %s from %s", e.ProjectRoleID, e.PersonID)
}

func (e *ProjectRoleUnassigned) Apply(project *project.Project) {
	for i, m := range project.Members {
		if m.PersonID == e.PersonID && m.ProjectRoleID == e.ProjectRoleID {
			project.Members = append(project.Members[:i], project.Members[i+1:]...)
			return
		}
	}
}

func (e *ProjectRoleUnassigned) RelatedIDs() RelatedIDs {
	return RelatedIDs{PersonID: &e.PersonID, ProjectRoleID: &e.ProjectRoleID}
}

func (e *ProjectRoleUnassigned) ApprovalMessage(projectTitle string) string {
	return fmt.Sprintf("Approval requested: role unassigned in project '%s'.", projectTitle)
}

type ProjectRoleUnassignedInput struct {
	PersonID      uuid.UUID `json:"person_id"`
	ProjectRoleID uuid.UUID `json:"project_role_id"`
}

func DecideProjectRoleUnassigned(
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *project.Project,
	in ProjectRoleUnassignedInput,
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

	found := false
	for _, m := range cur.Members {
		if m.PersonID == in.PersonID && m.ProjectRoleID == in.ProjectRoleID {
			found = true
			break
		}
	}
	if !found {
		return nil, nil
	}

	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = ProjectRoleUnassignedMeta.FriendlyName

	return &ProjectRoleUnassigned{
		Base:          base,
		PersonID:      in.PersonID,
		ProjectRoleID: in.ProjectRoleID,
	}, nil
}

var ProjectRoleUnassignedMeta = EventMeta{
	Type:         ProjectRoleUnassignedType,
	FriendlyName: "Project Role Unassignment",
}

func init() {
	RegisterMeta(ProjectRoleUnassignedMeta, func() Event {
		return &ProjectRoleUnassigned{
			Base: Base{FriendlyNameStr: ProjectRoleUnassignedMeta.FriendlyName},
		}
	})

	RegisterDecider[ProjectRoleUnassignedInput](ProjectRoleUnassignedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *project.Project, in ProjectRoleUnassignedInput, status Status) (Event, error) {
			return DecideProjectRoleUnassigned(projectID, actor, cur, in, status)
		})

	RegisterInputType(ProjectRoleUnassignedType, ProjectRoleUnassignedInput{})
}
