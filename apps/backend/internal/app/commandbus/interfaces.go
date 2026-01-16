package commandbus

import (
	"context"

	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
)

type EventStore interface {
	Load(ctx context.Context, id uuid.UUID) ([]events.Event, int, error)
	Append(ctx context.Context, id uuid.UUID, expectedVersion int, evts ...events.Event) error
}

type Publisher interface {
	HandleEvents(ctx context.Context, evts ...events.Event) error
}

type Reducer[T any] interface {
	Reduce(id uuid.UUID, history []events.Event) (*T, error)
	Apply(cur *T, e events.Event) error
}

type NewReducer[T any] interface {
	New(id uuid.UUID) *T
}
