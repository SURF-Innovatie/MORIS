package commandbus

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Executor[T any] struct {
	es  EventStore
	pub Publisher
	red Reducer[T]
	new NewReducer[T]
}

func NewExecutor[T any](es EventStore, pub Publisher, red Reducer[T], newOpt NewReducer[T]) *Executor[T] {
	return &Executor[T]{es: es, pub: pub, red: red, new: newOpt}
}

func (x *Executor[T]) Execute(ctx context.Context, id uuid.UUID, decide Decision[T]) (*T, error) {
	history, version, err := x.es.Load(ctx, id)
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

	for _, e := range newEvents {
		if err := x.es.Append(ctx, id, version, e); err != nil {
			return nil, fmt.Errorf("append %s: %w", e.Type(), err)
		}
		version++

		if err := x.red.Apply(cur, e); err != nil {
			return nil, fmt.Errorf("apply %s: %w", e.Type(), err)
		}
	}

	_ = x.pub.HandleEvents(ctx, newEvents...)

	return cur, nil
}
