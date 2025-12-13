package eventdto

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID `json:"id"`
	ProjectID    uuid.UUID `json:"projectId"`
	Type         string    `json:"type"`
	Status       string    `json:"status"`
	CreatedBy    uuid.UUID `json:"createdBy"`
	At           time.Time `json:"at"`
	Details      string    `json:"details"`
	ProjectTitle string    `json:"projectTitle"`

	// Optional “related object” pointers (IDs only)
	PersonID      *uuid.UUID `json:"personId,omitempty"`
	ProductID     *uuid.UUID `json:"productId,omitempty"`
	ProjectRoleID *uuid.UUID `json:"projectRoleId,omitempty"`
	OrgNodeID     *uuid.UUID `json:"orgNodeId,omitempty"`
}

type Response struct {
	Events []Event `json:"events"`
}

func FromEntity(e events.Event) Event {
	return FromEntityWithTitle(e, "")
}

func FromEntityWithTitle(e events.Event, projectTitle string) Event {
	createdBy := uuid.Nil
	if cb, ok := any(e).(interface{ CreatedByID() uuid.UUID }); ok {
		createdBy = cb.CreatedByID()
	}

	dtoEvent := Event{
		ID:           e.GetID(),
		ProjectID:    e.AggregateID(),
		Type:         e.Type(),
		Status:       e.GetStatus(),
		CreatedBy:    createdBy,
		At:           e.OccurredAt(),
		Details:      e.String(),
		ProjectTitle: projectTitle,
	}

	switch ev := e.(type) {
	case events.ProjectRoleAssigned:
		dtoEvent.PersonID = &ev.PersonID
		dtoEvent.ProjectRoleID = &ev.ProjectRoleID

	case events.ProjectRoleUnassigned:
		dtoEvent.PersonID = &ev.PersonID
		dtoEvent.ProjectRoleID = &ev.ProjectRoleID

	case events.ProductAdded:
		dtoEvent.ProductID = &ev.ProductID

	case events.ProductRemoved:
		dtoEvent.ProductID = &ev.ProductID

	case events.OwningOrgNodeChanged:
		dtoEvent.OrgNodeID = &ev.OwningOrgNodeID

	case events.ProjectStarted:
		dtoEvent.OrgNodeID = &ev.OwningOrgNodeID
	}

	return dtoEvent
}
