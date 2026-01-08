package dto

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID     `json:"id"`
	ProjectID    uuid.UUID     `json:"projectId"`
	Type         string        `json:"type"`
	Status       events.Status `json:"status"`
	CreatedBy    uuid.UUID     `json:"createdBy"`
	At           time.Time     `json:"at"`
	Details      string        `json:"details"`
	ProjectTitle string        `json:"projectTitle"`

	// Optional "related object" pointers (IDs only)
	PersonID      *uuid.UUID `json:"personId,omitempty"`
	ProductID     *uuid.UUID `json:"productId,omitempty"`
	ProjectRoleID *uuid.UUID `json:"projectRoleId,omitempty"`
	OrgNodeID     *uuid.UUID `json:"orgNodeId,omitempty"`

	// Optional full related objects
	Person      *PersonResponse      `json:"person,omitempty"`
	Product     *ProductResponse     `json:"product,omitempty"`
	ProjectRole *ProjectRoleResponse `json:"projectRole,omitempty"`
	// OrgNode     *OrganisationNodeResponse `json:"orgNode,omitempty"` // TODO: Add DTO if needed

	// The raw event data (input payload)
	Data any `json:"data,omitempty"`
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
	var coreEv events.Event = ev    
    return e.fromCore(coreEv, projectTitle)
}

func (e Event) FromDetailedEntity(dev events.DetailedEvent) Event {
	dto := e.fromCore(dev.Event, "")
	
	if dev.Person != nil {
		p := PersonResponse{}.FromEntity(*dev.Person)
		dto.Person = &p
	}
	if dev.Product != nil {
		p := ProductResponse{}.FromEntity(*dev.Product)
		dto.Product = &p
	}
	if dev.ProjectRole != nil {
		r := ProjectRoleResponse{
			ID: dev.ProjectRole.ID,
			Key: dev.ProjectRole.Key,
			Name: dev.ProjectRole.Name,
		}
		dto.ProjectRole = &r
	}
	
	return dto
}

func (e Event) fromCore(ev events.Event, projectTitle string) Event {
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
		Data:         ev,
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
