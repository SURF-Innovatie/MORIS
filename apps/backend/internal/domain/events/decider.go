package events

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Decision is a function that can emit 0..n events for a specific aggregate.
type Decision func(ctx context.Context, cur any, input any) ([]Event, error)

type ProjectDecider[I any] func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur any, input I, status Status) (Event, error)

// Registry: event type -> decider function
var eventDeciders = map[string]any{}

// RegisterDecider registers a decider for an event type.
func RegisterDecider[I any](eventType string, d ProjectDecider[I]) {
	eventDeciders[eventType] = d
}

func GetDecider(eventType string) (any, bool) {
	d, ok := eventDeciders[eventType]
	return d, ok
}

// ValidateRegistrations enforces: every event meta must have a decider.
func ValidateRegistrations() error {
	for typ := range eventMetas {
		if _, ok := eventDeciders[typ]; !ok {
			return fmt.Errorf("missing decider for event type %q", typ)
		}
	}
	return nil
}
