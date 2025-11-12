package domaintest

import (
	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/projection"
)

// Stream stores events per aggregate for tests.
type Stream map[uuid.UUID][]events.Event

func NewStream() Stream { return make(Stream) }

func (s Stream) Append(id uuid.UUID, ev events.Event) {
	s[id] = append(s[id], ev)
}

func (s Stream) Events(id uuid.UUID) []events.Event {
	return append([]events.Event(nil), s[id]...)
}

func (s Stream) Reduce(id uuid.UUID) *entities.Project {
	return projection.Reduce(id, s[id])
}
