package commandbus

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Executor[T any] struct {
	store EventStore
	pub   EventPublisher
	red   Reducer[T]
	new   NewReducer[T]
}

func NewExecutor[T any](store EventStore, pub EventPublisher, red Reducer[T], newOpt NewReducer[T]) *Executor[T] {
	return &Executor[T]{store: store, pub: pub, red: red, new: newOpt}
}

func (x *Executor[T]) Execute(ctx context.Context, id uuid.UUID, decide Decision[T]) (*T, error) {
	history, version, err := x.store.Load(ctx, id)
	if err != nil {
		return nil, err
	}

	var cur *T
	if len(history) == 0 {
		if x.new == nil {
			return nil, fmt.Errorf("aggregate %s not found", id)
		}
		cur = x.new.New(id)
		version = 0
	} else {
		cur, err = x.red.Reduce(id, history)
		if err != nil {
			return nil, err
		}
	}

	newEvents, err := decide(ctx, cur)
	if err != nil {
		return nil, err
	}
	if len(newEvents) == 0 {
		return cur, nil
	}

	// Prefer ONE append with expectedVersion for the whole batch if your store supports it.
	// If not, keep the loop; but note that optimistic concurrency is typically per-stream.
	if err := x.store.Append(ctx, id, version, newEvents...); err != nil {
		return nil, fmt.Errorf("append: %w", err)
	}

	// Update in-memory state
	for _, e := range newEvents {
		if err := x.red.Apply(cur, e); err != nil {
			return nil, fmt.Errorf("apply %s: %w", e.Type(), err)
		}
	}

	// Publish side effects (infra dispatcher). No longer HandleEvents on event service.
	if x.pub != nil {
		_ = x.pub.Publish(ctx, newEvents...)
	}

	return cur, nil
}
