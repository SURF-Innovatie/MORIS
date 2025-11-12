package events

import (
	"time"

	"github.com/google/uuid"
)

// Event is a sealed interface for all domain events.
type Event interface {
	isEvent()
	// AggregateID identifies the Project this event belongs to.
	AggregateID() uuid.UUID
	// OccurredAt is when the change happened (UTC).
	OccurredAt() time.Time
	// Type is a stable string for routing/serialization.
	Type() string
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

// Base carries common metadata. Embed in all events.
type Base struct {
	ProjectID uuid.UUID `json:"projectId"`
	At        time.Time `json:"at"`
}

func (b Base) AggregateID() uuid.UUID { return b.ProjectID }
func (b Base) OccurredAt() time.Time  { return b.At }
