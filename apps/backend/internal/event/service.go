package event

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/app/notification"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	ApproveEvent(ctx context.Context, eventID uuid.UUID) error
	RejectEvent(ctx context.Context, eventID uuid.UUID) error
	HandleEvents(ctx context.Context, evts ...events.Event) error
	RegisterNotificationHandler(handler NotificationHandler)
	RegisterStatusChangeHandler(handler StatusChangeHandler)
	GetEvent(ctx context.Context, eventID uuid.UUID) (events.Event, error)
	GetEventTypes(ctx context.Context) ([]events.EventMeta, error)
}

type StatusChangeHandler func(ctx context.Context, event events.Event) error

type NotificationHandler interface {
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
			if err := handler.Handle(ctx, e); err != nil {
				log.Error().Err(err).Msgf("Error handling event %s with handler %T", e.GetID(), handler)
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

	// 2. Mark related notifications as read
	if err := s.notifier.MarkAsReadByEventID(ctx, eventID); err != nil {
		log.Warn().Err(err).Msgf("Failed to mark notifications as read for event %s", eventID)
	}

	event, err := s.es.LoadEvent(ctx, eventID)
	if err != nil {
		return err
	}

	// Notify status change handlers
	for _, h := range s.statusChangeHandlers {
		if err := h(ctx, event); err != nil {
			log.Error().Err(err).Msg("Error in status change handler")
		}
	}

	return s.HandleEvents(ctx, event)
}

func (s *service) RejectEvent(ctx context.Context, eventID uuid.UUID) error {
	if err := s.es.UpdateEventStatus(ctx, eventID, "rejected"); err != nil {
		return err
	}

	// Mark related notifications as read
	if err := s.notifier.MarkAsReadByEventID(ctx, eventID); err != nil {
		log.Warn().Err(err).Msgf("Failed to mark notifications as read for event %s", eventID)
	}

	event, err := s.es.LoadEvent(ctx, eventID)
	if err != nil {
		return err
	}

	// Notify status change handlers
	for _, h := range s.statusChangeHandlers {
		if err := h(ctx, event); err != nil {
			log.Error().Err(err).Msg("Error in status change handler")
		}
	}

	return s.HandleEvents(ctx, event)
}

func (s *service) GetEvent(ctx context.Context, eventID uuid.UUID) (events.Event, error) {
	return s.es.LoadEvent(ctx, eventID)
}

func (s *service) GetEventTypes(ctx context.Context) ([]events.EventMeta, error) {
	metas := events.GetAllMetas()

	return metas, nil
}
