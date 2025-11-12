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
		p.Organisation = ev.Organisation
		p.People = toPersonPtrSlice(ev.People)

	case events.TitleChanged:
		p.Title = ev.Title

	case events.DescriptionChanged:
		p.Description = ev.Description

	case events.StartDateChanged:
		p.StartDate = ev.StartDate

	case events.EndDateChanged:
		p.EndDate = ev.EndDate

	case events.OrganisationChanged:
		p.Organisation = ev.Organisation

	case events.PersonAdded:
		if !hasPerson(p.People, ev.Person.Name) {
			p.People = append(p.People, &entities.Person{Name: ev.Person.Name})
		}

	case events.PersonRemoved:
		p.People = filterPeople(p.People, func(pe *entities.Person) bool {
			return pe.Name != ev.Name
		})

	default:
		// unknown event type: ignore to keep forward-compatible
	}
}

func toPersonPtrSlice(in []entities.Person) []*entities.Person {
	if len(in) == 0 {
		return nil
	}
	out := make([]*entities.Person, 0, len(in))
	for i := range in {
		pe := in[i] // take address of loop variable copy
		out = append(out, &pe)
	}
	return out
}

func hasPerson(list []*entities.Person, name string) bool {
	for _, p := range list {
		if p != nil && p.Name == name {
			return true
		}
	}
	return false
}

func filterPeople(list []*entities.Person, keep func(*entities.Person) bool) []*entities.Person {
	if len(list) == 0 {
		return nil
	}
	out := list[:0]
	for _, p := range list {
		if p != nil && keep(p) {
			out = append(out, p)
		}
	}
	return out
}
