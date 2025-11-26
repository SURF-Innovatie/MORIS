package events

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	isEvent()
	AggregateID() uuid.UUID
	OccurredAt() time.Time
	Type() string
	String() string
}

const (
	ProjectStartedType      = "project.started"
	TitleChangedType        = "project.title_changed"
	DescriptionChangedType  = "project.description_changed"
	StartDateChangedType    = "project.start_date_changed"
	EndDateChangedType      = "project.end_date_changed"
	OrganisationChangedType = "project.organisation_changed"
	PersonAddedType         = "project.person_added"
	PersonRemovedType       = "project.person_removed"
)

type Base struct {
	ProjectID uuid.UUID `json:"projectId"`
	At        time.Time `json:"at"`
}

func (b Base) AggregateID() uuid.UUID { return b.ProjectID }
func (b Base) OccurredAt() time.Time  { return b.At }
