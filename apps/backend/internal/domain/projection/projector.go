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
	switch ev := e.(type) {
	case events.ProjectStarted:
		p.Title = ev.Title
		p.Description = ev.Description
		p.StartDate = ev.StartDate
		p.EndDate = ev.EndDate
		p.Organisation = ev.OrganisationID
		p.People = ev.People

	case events.TitleChanged:
		p.Title = ev.Title

	case events.DescriptionChanged:
		p.Description = ev.Description

	case events.StartDateChanged:
		p.StartDate = ev.StartDate

	case events.EndDateChanged:
		p.EndDate = ev.EndDate

	case events.OrganisationChanged:
		p.Organisation = ev.OrganisationID

	case events.PersonAdded:
		if !hasPerson(p.People, ev.PersonId) {
			p.People = append(p.People, ev.PersonId)
		}

	case events.PersonRemoved:
		p.People = filterPeople(p.People, func(id uuid.UUID) bool {
			return id != ev.PersonId
		})

	default:
		// unknown event type: ignore to keep forward-compatible
	}
}

func hasPerson(list []uuid.UUID, id uuid.UUID) bool {
	for _, p := range list {
		if p == id {
			return true
		}
	}
	return false
}

func filterPeople(list []uuid.UUID, keep func(uuid.UUID) bool) []uuid.UUID {
	if len(list) == 0 {
		return nil
	}
	out := list[:0]
	for _, p := range list {
		if keep(p) {
			out = append(out, p)
		}
	}
	return out
}
