package eventdispatch

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/rs/zerolog/log"
)

type NotificationHandler interface {
	Handle(ctx context.Context, event events.Event) error
}

type StatusChangeHandler func(ctx context.Context, event events.Event) error

type Dispatcher struct {
	notificationHandlers []NotificationHandler
	statusChangeHandlers []StatusChangeHandler
}

func New(
	notificationHandlers []NotificationHandler,
	statusChangeHandlers []StatusChangeHandler,
) *Dispatcher {
	return &Dispatcher{
		notificationHandlers: notificationHandlers,
		statusChangeHandlers: statusChangeHandlers,
	}
}

func (d *Dispatcher) Publish(ctx context.Context, evts ...events.Event) error {
	for _, e := range evts {
		for _, h := range d.notificationHandlers {
			if err := h.Handle(ctx, e); err != nil {
				log.Error().Err(err).Msgf("Error handling event %s with handler %T", e.GetID(), h)
			}
		}
	}
	return nil
}

func (d *Dispatcher) PublishStatusChanged(ctx context.Context, evt events.Event) error {
	for _, h := range d.statusChangeHandlers {
		if err := h(ctx, evt); err != nil {
			log.Error().Err(err).Msg("Error in status change handler")
		}
	}
	return nil
}
