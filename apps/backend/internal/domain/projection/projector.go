package projection

import (
	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
)

// Reduce rehydrates a Project from its event stream.
func Reduce(id uuid.UUID, es []events.Event) *entities.Project {
	p := &entities.Project{Id: id}
	for _, e := range es {
		Apply(p, e)
	}
	return p
}

// Apply mutates the given Project based on one event.
func Apply(p *entities.Project, e events.Event) {
	if e.GetStatus() == "pending" || e.GetStatus() == "rejected" {
		return
	}

	switch ev := e.(type) {
	case events.ProjectStarted:
		p.Title = ev.Title
		p.Description = ev.Description
		p.StartDate = ev.StartDate
		p.EndDate = ev.EndDate
		p.OwningOrgNodeID = ev.OwningOrgNodeID
		p.Members = ev.Members

	case events.OwningOrgNodeChanged:
		p.OwningOrgNodeID = ev.OwningOrgNodeID

	case events.ProjectRoleAssigned:
		for _, m := range p.Members {
			if m.PersonID == ev.PersonID && m.ProjectRoleID == ev.ProjectRoleID {
				return // already assigned
			}
		}
		p.Members = append(p.Members, entities.ProjectMember{
			PersonID:      ev.PersonID,
			ProjectRoleID: ev.ProjectRoleID,
		})

	case events.ProjectRoleUnassigned:
		shouldRemove := -1
		for i, m := range p.Members {
			if m.PersonID == ev.PersonID && m.ProjectRoleID == ev.ProjectRoleID {
				shouldRemove = i
				break
			}
		}
		if shouldRemove != -1 {
			p.Members = append(p.Members[:shouldRemove], p.Members[shouldRemove+1:]...)
		}

	case events.TitleChanged:
		p.Title = ev.Title

	case events.DescriptionChanged:
		p.Description = ev.Description

	case events.StartDateChanged:
		p.StartDate = ev.StartDate

	case events.EndDateChanged:
		p.EndDate = ev.EndDate

	case events.ProductAdded:
		for _, prod := range p.ProductIDs {
			if prod == ev.ProductID {
				return
			}
		}
		p.ProductIDs = append(p.ProductIDs, ev.ProductID)

	case events.ProductRemoved:
		shouldRemove := -1
		for i, p := range p.ProductIDs {
			if p == ev.ProductID {
				shouldRemove = i
				break
			}
		}
		if shouldRemove != -1 {
			p.ProductIDs = append(p.ProductIDs[:shouldRemove], p.ProductIDs[shouldRemove+1:]...)
		}

	default:
		// unknown event type: ignore to keep forward-compatible
	}
}
