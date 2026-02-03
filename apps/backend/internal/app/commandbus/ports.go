package commandbus

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/google/uuid"
)

type EventStore interface {
	Load(ctx context.Context, id uuid.UUID) ([]events.Event, int, error)
	Append(ctx context.Context, id uuid.UUID, expectedVersion int, evts ...events.Event) error
}

type EventPublisher interface {
	Publish(ctx context.Context, evts ...events.Event) error
}

type Reducer[T any] interface {
	Reduce(id uuid.UUID, history []events.Event) (*T, error)
	Apply(cur *T, e events.Event) error
}

type NewReducer[T any] interface {
	New(id uuid.UUID) *T
}
