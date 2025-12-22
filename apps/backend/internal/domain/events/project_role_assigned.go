package events

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

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

func init() {
	RegisterMeta(EventMeta{
		Type:         ProjectRoleAssignedType,
		FriendlyName: "Project Role Assignment",
		CheckApproval: func(ctx context.Context, event Event, client *ent.Client) bool {
			return true // Always requires approval
		},
	}, func() Event { return &ProjectRoleAssigned{} })
}
