package eventdto

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/api/persondto"
	"github.com/SURF-Innovatie/MORIS/internal/api/productdto"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID            `json:"id"`
	ProjectID    uuid.UUID            `json:"projectId"`
	Type         string               `json:"type"`
	Status       string               `json:"status"`
	CreatedBy    uuid.UUID            `json:"createdBy"`
	At           time.Time            `json:"at"`
	Details      string               `json:"details"`      // Human readable description
	ProjectTitle string               `json:"projectTitle"` // Title of the project
	Person       *persondto.Response  `json:"person,omitempty"`
	Product      *productdto.Response `json:"product,omitempty"`
}

type Response struct {
	Events []Event `json:"events"`
}

func FromEntity(e events.Event) Event {
	return FromEntityWithTitle(e, "")
}

func FromEntityWithTitle(e events.Event, projectTitle string) Event {
	dtoEvent := Event{
		ID:           e.GetID(),
		ProjectID:    e.AggregateID(),
		Type:         e.Type(),
		Status:       e.GetStatus(),
		CreatedBy:    e.(interface{ CreatedByID() uuid.UUID }).CreatedByID(), // Assuming CreatedByID exists on all events or checking interface
		At:           e.OccurredAt(),
		Details:      e.String(),
		ProjectTitle: projectTitle,
	}

	switch ev := e.(type) {
	case events.PersonAdded:
		p := persondto.FromEntity(ev.Person)
		dtoEvent.Person = &p
	case events.PersonRemoved:
		p := persondto.FromEntity(ev.Person)
		dtoEvent.Person = &p
	case events.ProductAdded:
		p := productdto.FromEntity(ev.Product)
		dtoEvent.Product = &p
	case events.ProductRemoved:
		p := productdto.FromEntity(ev.Product)
		dtoEvent.Product = &p
	}

	return dtoEvent
}
