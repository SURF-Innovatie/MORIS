package project

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/identity"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/role"
	"github.com/google/uuid"
)

type Member struct {
	PersonID      uuid.UUID
	ProjectRoleID uuid.UUID
}

type Project struct {
	Id              uuid.UUID
	Version         int
	StartDate       time.Time
	EndDate         time.Time
	Title           string
	Description     string
	Members         []Member
	OwningOrgNodeID uuid.UUID
	ProductIDs      []uuid.UUID
	CustomFields    map[string]interface{}
}

type MemberDetail struct {
	Person identity.Person
	Role   role.ProjectRole
}

type ChangeLogEntry struct {
	Event string
	At    time.Time
}

type ChangeLog struct {
	Entries []ChangeLogEntry
}
