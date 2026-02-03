package projection

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/google/uuid"
)

// Reduce rehydrates a Project from its event stream.
func Reduce(id uuid.UUID, es []events.Event) *project.Project {
	p := &project.Project{Id: id}
	for _, e := range es {
		Apply(p, e)
	}
	return p
}

// Apply mutates the given Project based on one event.
func Apply(p *project.Project, e events.Event) {
	if p == nil {
		return
	}

	// TODO: Rework this when implementing event statuses properly
	if e.GetStatus() == "pending" || e.GetStatus() == "rejected" {
		return
	}

	if applier, ok := e.(events.Applier); ok {
		applier.Apply(p)
	}
}
