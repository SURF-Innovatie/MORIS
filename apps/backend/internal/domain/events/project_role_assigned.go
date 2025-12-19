package events

import (
	"fmt"

	"github.com/google/uuid"
)

type ProjectRoleAssigned struct {
	Base
	PersonID       uuid.UUID `json:"person_id"`
	ProjectRoleKey string    `json:"project_role_key"`
}

func (ProjectRoleAssigned) isEvent()     {}
func (ProjectRoleAssigned) Type() string { return ProjectRoleAssignedType }
func (e ProjectRoleAssigned) String() string {
	return fmt.Sprintf("Project role assigned: %s to %s", e.ProjectRoleKey, e.PersonID)
}
