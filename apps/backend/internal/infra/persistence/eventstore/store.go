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

	// UpdateEventStatus updates the status of an event.
	UpdateEventStatus(ctx context.Context, eventID uuid.UUID, status string) error

	// LoadEvent loads a single event by ID.
	LoadEvent(ctx context.Context, eventID uuid.UUID) (events.Event, error)

	// LoadUserApprovedEvents loads all approved events created by a user.
	LoadUserApprovedEvents(ctx context.Context, userID uuid.UUID) ([]events.Event, error)
}

var ErrConcurrency = errors.New("concurrency conflict")
