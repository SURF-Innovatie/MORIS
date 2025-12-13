package events

import (
	"fmt"

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
	return fmt.Sprintf("Project role assigned: %s to %s", e.ProjectRoleID, e.PersonID)
}
