package events

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	isEvent()
	AggregateID() uuid.UUID
	OccurredAt() time.Time
	GetID() uuid.UUID
	Type() string
	String() string
	CreatedByID() uuid.UUID
	GetStatus() string
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
	ProductAddedType        = "project.product_added"
	ProductRemovedType      = "project.product_removed"
)

type Base struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"projectId"`
	At        time.Time `json:"at"`
	CreatedBy uuid.UUID `json:"createdBy"`
	Status    string    `json:"status"`
}

func (b Base) AggregateID() uuid.UUID { return b.ProjectID }
func (b Base) OccurredAt() time.Time  { return b.At }
func (b Base) GetID() uuid.UUID       { return b.ID }
func (b Base) CreatedByID() uuid.UUID { return b.CreatedBy }
func (b Base) GetStatus() string      { return b.Status }
