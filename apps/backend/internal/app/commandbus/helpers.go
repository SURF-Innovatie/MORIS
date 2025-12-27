package commandbus

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
)

type Decision[T any] func(ctx context.Context, cur *T) ([]events.Event, error)

func One(e events.Event, err error) ([]events.Event, error) {
	if err != nil || e == nil {
		return nil, err
	}
	return []events.Event{e}, nil
}

func Add(out *[]events.Event, e events.Event, err error) error {
	if err != nil {
		return err
	}
	if e != nil {
		*out = append(*out, e)
	}
	return nil
}
