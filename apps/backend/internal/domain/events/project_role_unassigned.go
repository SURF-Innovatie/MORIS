package events

import (
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
	return "Project role unassigned: " + e.ProjectRoleID.String() + " from " + e.PersonID.String()
}
