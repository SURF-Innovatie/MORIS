package dto

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

	// Optional "related object" pointers (IDs only)
	PersonID      *uuid.UUID `json:"personId,omitempty"`
	ProductID     *uuid.UUID `json:"productId,omitempty"`
	ProjectRoleID *uuid.UUID `json:"projectRoleId,omitempty"`
	OrgNodeID     *uuid.UUID `json:"orgNodeId,omitempty"`
}

type EventResponse struct {
	Events []Event `json:"events"`
}

type EventTypeResponse struct {
	Type         string `json:"type"`
	FriendlyName string `json:"friendlyName"`
	Allowed      bool   `json:"allowed"`
	Description  string `json:"description,omitempty"`
}

func (e Event) FromEntity(ev events.Event) Event {
	return e.FromEntityWithTitle(ev, "")
}

func (e Event) FromEntityWithTitle(ev events.Event, projectTitle string) Event {
	createdBy := uuid.Nil
	if cb, ok := any(ev).(interface{ CreatedByID() uuid.UUID }); ok {
		createdBy = cb.CreatedByID()
	}

	dtoEvent := Event{
		ID:           ev.GetID(),
		ProjectID:    ev.AggregateID(),
		Type:         ev.Type(),
		Status:       ev.GetStatus(),
		CreatedBy:    createdBy,
		At:           ev.OccurredAt(),
		Details:      ev.String(),
		ProjectTitle: projectTitle,
	}

	// Enrich with related IDs if available
	if r, ok := ev.(events.HasRelatedIDs); ok {
		ids := r.RelatedIDs()
		dtoEvent.PersonID = ids.PersonID
		dtoEvent.ProductID = ids.ProductID
		dtoEvent.ProjectRoleID = ids.ProjectRoleID
		dtoEvent.OrgNodeID = ids.OrgNodeID
	}

	return dtoEvent
}
