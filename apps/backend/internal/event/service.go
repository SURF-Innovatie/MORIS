package event

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	"github.com/SURF-Innovatie/MORIS/internal/notification"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Service interface {
	ApproveEvent(ctx context.Context, eventID uuid.UUID) error
	RejectEvent(ctx context.Context, eventID uuid.UUID) error
	HandleEvents(ctx context.Context, evts ...events.Event) error
	RegisterNotificationHandler(handler NotificationHandler)
	RegisterStatusChangeHandler(handler StatusChangeHandler)
	GetEvent(ctx context.Context, eventID uuid.UUID) (events.Event, error)
}

type StatusChangeHandler func(ctx context.Context, event events.Event) error

type NotificationHandler interface {
	CanHandle(event events.Event) bool
	Handle(ctx context.Context, event events.Event) error
}

type service struct {
	cli                  *ent.Client
	es                   eventstore.Store
	notifier             notification.Service
	notificationHandlers []NotificationHandler
	statusChangeHandlers []StatusChangeHandler
}

func NewService(es eventstore.Store, cli *ent.Client, notifier notification.Service) Service {
	return &service{es: es, cli: cli, notifier: notifier}
}

func (s *service) RegisterNotificationHandler(handler NotificationHandler) {
	s.notificationHandlers = append(s.notificationHandlers, handler)
}

func (s *service) RegisterStatusChangeHandler(handler StatusChangeHandler) {
	s.statusChangeHandlers = append(s.statusChangeHandlers, handler)
}

func (s *service) HandleEvents(ctx context.Context, evts ...events.Event) error {
	for _, e := range evts {
		// Existing logic integrated via handlers now
		for _, handler := range s.notificationHandlers {
			if handler.CanHandle(e) {
				if err := handler.Handle(ctx, e); err != nil {
					logrus.Errorf("Error handling event %s with handler %T: %v\n", e.GetID(), handler, err)
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

	event, err := s.es.LoadEvent(ctx, eventID)
	if err != nil {
		return err
	}

	// Notify status change handlers
	for _, h := range s.statusChangeHandlers {
		if err := h(ctx, event); err != nil {
			logrus.Errorf("Error in status change handler: %v\n", err)
		}
	}

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

	// Notify status change handlers
	for _, h := range s.statusChangeHandlers {
		if err := h(ctx, event); err != nil {
			logrus.Errorf("Error in status change handler: %v\n", err)
		}
	}

	return s.HandleEvents(ctx, event)
}

func (s *service) GetEvent(ctx context.Context, eventID uuid.UUID) (events.Event, error) {
	return s.es.LoadEvent(ctx, eventID)
}
