package adapter

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
)

// SourceAdapter imports data from external systems as domain events.
type SourceAdapter interface {
	// Name returns the unique identifier for this adapter.
	Name() string

	// DisplayName returns the human-readable name for this adapter.
	DisplayName() string

	// SupportedTypes returns which data types this source can provide.
	SupportedTypes() []DataType

	// InputInfo returns information about the input fields required for this adapter.
	InputInfo() InputInfo

	// Connect establishes connection to the external system.
	Connect(ctx context.Context) error

	// FetchProjects streams ProjectStarted events for import.
	// Returns a channel of events and a channel for errors.
	FetchProjects(ctx context.Context, opts FetchOptions) (<-chan events.Event, <-chan error)

	// FetchUsers streams user/person data for import.
	// Returns a channel of UserContext and a channel for errors.
	FetchUsers(ctx context.Context, opts FetchOptions) (<-chan *UserContext, <-chan error)

	// Close releases any resources.
	Close() error
}
