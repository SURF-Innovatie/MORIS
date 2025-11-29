package entities

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	Id           uuid.UUID
	ProjectAdmin uuid.UUID
	Version      int
	StartDate    time.Time
	EndDate      time.Time
	Title        string
	Description  string
	People       []uuid.UUID
	Organisation uuid.UUID
	Products     []uuid.UUID
}

type ProjectDetails struct {
	Project      Project
	Organisation Organisation
	People       []Person
	Products     []Product
}

type ChangeLogEntry struct {
	Event string
	At    time.Time
}

type ChangeLog struct {
	Entries []ChangeLogEntry
}
