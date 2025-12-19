package events

import (
	"github.com/google/uuid"
)

type ProjectRoleUnassigned struct {
	Base
	PersonID       uuid.UUID `json:"person_id"`
	ProjectRoleKey string    `json:"project_role_key"`
}

func (ProjectRoleUnassigned) isEvent()     {}
func (ProjectRoleUnassigned) Type() string { return ProjectRoleUnassignedType }
func (e ProjectRoleUnassigned) String() string {
	return "Project role unassigned: " + e.ProjectRoleKey + " from " + e.PersonID.String()
}
