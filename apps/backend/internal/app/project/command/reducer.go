package command

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/projection"
	"github.com/google/uuid"
)

type Reducer struct{}

func (Reducer) Reduce(id uuid.UUID, history []events.Event) (*project.Project, error) {
	p := projection.Reduce(id, history)
	return p, nil
}

func (Reducer) Apply(cur *project.Project, e events.Event) error {
	projection.Apply(cur, e)
	return nil
}

type NewReducer struct{}

func (NewReducer) New(id uuid.UUID) *project.Project {
	return &project.Project{Id: id}
}
