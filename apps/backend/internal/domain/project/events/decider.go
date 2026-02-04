package events

import (
	"context"
	"encoding/json"
	"fmt"

	projdomain "github.com/SURF-Innovatie/MORIS/internal/domain/project"
	"github.com/google/uuid"
)

type ProjectDecider[I any] func(
	ctx context.Context,
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *projdomain.Project,
	input I,
	status Status,
) (Event, error)

type ProjectDeciderRaw func(
	ctx context.Context,
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *projdomain.Project,
	input json.RawMessage,
	status Status,
) (Event, error)

var eventDeciders = map[string]ProjectDeciderRaw{}

// RegisterDecider keeps typed ergonomics but stores a raw-json wrapper.
func RegisterDecider[I any](eventType string, d ProjectDecider[I]) {
	eventDeciders[eventType] = func(
		ctx context.Context,
		projectID uuid.UUID,
		actor uuid.UUID,
		cur *projdomain.Project,
		raw json.RawMessage,
		status Status,
	) (Event, error) {
		var in I
		if err := json.Unmarshal(raw, &in); err != nil {
			return nil, err
		}
		return d(ctx, projectID, actor, cur, in, status)
	}
}

func GetDecider(eventType string) (ProjectDeciderRaw, bool) {
	d, ok := eventDeciders[eventType]
	return d, ok
}

func ValidateRegistrations() error {
	for typ := range eventMetas {
		if _, ok := eventDeciders[typ]; !ok {
			return fmt.Errorf("missing decider for event type %q", typ)
		}
	}
	return nil
}
