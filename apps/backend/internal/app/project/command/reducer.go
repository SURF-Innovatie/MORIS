package command

import (
	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/projection"
)

type Reducer struct{}

func (Reducer) Reduce(id uuid.UUID, history []events.Event) (*entities.Project, error) {
	p := projection.Reduce(id, history)
	return p, nil
}

func (Reducer) Apply(cur *entities.Project, e events.Event) error {
	projection.Apply(cur, e)
	return nil
}

type NewReducer struct{}

func (NewReducer) New(id uuid.UUID) *entities.Project {
	return &entities.Project{Id: id}
}
