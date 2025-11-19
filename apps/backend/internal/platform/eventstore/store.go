package eventstore

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
)

type Store interface {
	// Load returns all events for this project and the current version.
	Load(ctx context.Context, id uuid.UUID) ([]events.Event, int, error)

	// Append appends newEvents, assuming the current version is expectedVersion.
	// Should return ErrConcurrency if the version is not as expected.
	Append(ctx context.Context, id uuid.UUID, expectedVersion int, newEvents ...events.Event) error
}

var ErrConcurrency = errors.New("concurrency conflict")
