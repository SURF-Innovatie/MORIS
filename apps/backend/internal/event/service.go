package event

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	"github.com/SURF-Innovatie/MORIS/internal/notification"
	"github.com/google/uuid"
)

type Service interface {
	ApproveEvent(ctx context.Context, eventID uuid.UUID) error
	RejectEvent(ctx context.Context, eventID uuid.UUID) error
	HandleEvents(ctx context.Context, evts ...events.Event) error
	RegisterNotificationHandler(handler NotificationHandler)
	GetEvent(ctx context.Context, eventID uuid.UUID) (events.Event, error)
}

type NotificationHandler interface {
	CanHandle(event events.Event) bool
	Handle(ctx context.Context, event events.Event) error
}

type service struct {
	cli                  *ent.Client
	es                   eventstore.Store
	notifier             notification.Service
	notificationHandlers []NotificationHandler
}

func NewService(es eventstore.Store, cli *ent.Client, notifier notification.Service) Service {
	return &service{es: es, cli: cli, notifier: notifier}
}

func (s *service) RegisterNotificationHandler(handler NotificationHandler) {
	s.notificationHandlers = append(s.notificationHandlers, handler)
}

func (s *service) HandleEvents(ctx context.Context, evts ...events.Event) error {
	for _, e := range evts {
		// Existing logic integrated via handlers now
		for _, handler := range s.notificationHandlers {
			if handler.CanHandle(e) {
				if err := handler.Handle(ctx, e); err != nil {
					// Log error but continue?
					// Ideally we don't want to stop processing other handlers or events
					// use logrus if available or fmt
					fmt.Printf("Error handling event %s with handler %T: %v\n", e.GetID(), handler, err)
				}
			}
		}
	}
	return nil
}

func (s *service) ApproveEvent(ctx context.Context, eventID uuid.UUID) error {
	// 1. Update status
	if err := s.es.UpdateEventStatus(ctx, eventID, "approved"); err != nil {
		return err
	}

	// 2. Notify creator - Now handled via HandleEvents if we trigger it
	// But wait, UpdateEventStatus updates the event in the store but doesn't return it as a "new" event to process?
	// The `HandleEvents` expects `events.Event` objects.
	// We need to load the event to pass it to HandleEvents.

	event, err := s.es.LoadEvent(ctx, eventID)
	if err != nil {
		return err
	}

	// We need to manually trigger HandleEvents because this is a status change NOT originating from a new event append in the same way?
	// Or should we construct a wrapping event?
	// The `NotifyStatusUpdate` handler expects an event with status "approved".
	// The event loaded from store SHOULD have the new status "approved" because we just updated it.

	return s.HandleEvents(ctx, event)
}

func (s *service) RejectEvent(ctx context.Context, eventID uuid.UUID) error {
	if err := s.es.UpdateEventStatus(ctx, eventID, "rejected"); err != nil {
		return err
	}

	event, err := s.es.LoadEvent(ctx, eventID)
	if err != nil {
		return err
	}

	return s.HandleEvents(ctx, event)
}

func (s *service) GetEvent(ctx context.Context, eventID uuid.UUID) (events.Event, error) {
	return s.es.LoadEvent(ctx, eventID)
}
