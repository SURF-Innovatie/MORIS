package dto

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	events2 "github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID      `json:"id"`
	ProjectID    uuid.UUID      `json:"projectId"`
	Type         string         `json:"type"`
	Status       events2.Status `json:"status"`
	CreatedBy    uuid.UUID      `json:"createdBy"`
	At           time.Time      `json:"at"`
	Details      string         `json:"details"`
	ProjectTitle string         `json:"projectTitle"`
	FriendlyName string         `json:"friendlyName,omitempty"`

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

	Creator *PersonResponse `json:"creator,omitempty"`

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

func (e Event) FromEntity(ev events2.Event) Event {
	return e.FromEntityWithTitle(ev, "")
}

func (e Event) FromEntityWithTitle(ev events2.Event, projectTitle string) Event {
	var coreEv = ev
	return e.fromCore(coreEv, projectTitle)
}

func (e Event) FromDetailedEntity(dev events2.DetailedEvent) Event {
	dto := e.fromCore(dev.Event, "")

	if dev.Person != nil {
		p := transform.ToDTOItem[PersonResponse](*dev.Person)
		dto.Person = &p
	}
	if dev.Product != nil {
		p := transform.ToDTOItem[ProductResponse](*dev.Product)
		dto.Product = &p
	}
	if dev.ProjectRole != nil {
		r := transform.ToDTOItem[ProjectRoleResponse](*dev.ProjectRole)
		dto.ProjectRole = &r
	}
	if dev.Creator != nil {
		p := transform.ToDTOItem[PersonResponse](*dev.Creator)
		dto.Creator = &p
	}

	return dto
}

func (e Event) fromCore(ev events2.Event, projectTitle string) Event {
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
		FriendlyName: ev.FriendlyName(),
		Data:         ev,
	}

	// Enrich with related IDs if available
	if r, ok := ev.(events2.HasRelatedIDs); ok {
		ids := r.RelatedIDs()
		dtoEvent.PersonID = ids.PersonID
		dtoEvent.ProductID = ids.ProductID
		dtoEvent.ProjectRoleID = ids.ProjectRoleID
		dtoEvent.OrgNodeID = ids.OrgNodeID
	}

	return dtoEvent
}
