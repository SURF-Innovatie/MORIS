package adapter

import "context"

// SinkAdapter exports data to external systems.
type SinkAdapter interface {
	// Name returns the unique identifier for this adapter.
	Name() string

	// DisplayName returns the human-readable name for this adapter.
	DisplayName() string

	// SupportedTypes returns which data types this sink accepts.
	SupportedTypes() []DataType
	// OutputInfo returns information about the output destination of this adapter.
	OutputInfo() OutputInfo

	// Connect establishes connection to the external system.
	Connect(ctx context.Context) error

	// PushProject exports a project (with its event stream) to external system.
	PushProject(ctx context.Context, project ProjectContext) error

	// PushUser exports a user to external system.
	PushUser(ctx context.Context, user UserContext) error

	// Close releases any resources.
	Close() error
}
