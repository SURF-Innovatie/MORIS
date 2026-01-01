package entities

import (
	"time"

	"github.com/google/uuid"
)

type ProjectMember struct {
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
	Members         []ProjectMember
	OwningOrgNodeID uuid.UUID
	ProductIDs      []uuid.UUID
}

type ProjectMemberDetail struct {
	Person Person
	Role   ProjectRole
}

type ChangeLogEntry struct {
	Event string
	At    time.Time
}

type ChangeLog struct {
	Entries []ChangeLogEntry
}
