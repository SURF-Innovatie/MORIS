package entities

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	Id           uuid.UUID
	Version      int
	StartDate    time.Time
	EndDate      time.Time
	Title        string
	Description  string
	People       []*Person
	Organisation Organisation
}
