package events

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

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

func (e *ProjectRoleUnassigned) Apply(project *entities.Project) {
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

func init() {
	RegisterMeta(EventMeta{
		Type:         ProjectRoleUnassignedType,
		FriendlyName: "Project Role Unassignment",
		CheckApproval: func(ctx context.Context, event Event, client *ent.Client) bool {
			return true // Always requires approval
		},
	}, func() Event { return &ProjectRoleUnassigned{} })
}
