package event

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/app/notification"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	ApproveEvent(ctx context.Context, eventID uuid.UUID) error
	RejectEvent(ctx context.Context, eventID uuid.UUID) error
	GetEvent(ctx context.Context, eventID uuid.UUID) (events.Event, error)
	GetEventTypes(ctx context.Context) ([]events.EventMeta, error)
	LoadUserApprovedEvents(ctx context.Context, userID uuid.UUID) ([]events.Event, error)
	Load(ctx context.Context, id uuid.UUID) ([]events.Event, int, error)
	Append(ctx context.Context, id uuid.UUID, expectedVersion int, newEvents ...events.Event) error
	UpdateStatus(ctx context.Context, eventID uuid.UUID, status string) error
}

type StatusChangeHandler func(ctx context.Context, event events.Event) error

type NotificationHandler interface {
	Handle(ctx context.Context, event events.Event) error
}

type service struct {
	repo      repository
	notifier  notification.Service
	publisher Publisher
}

func NewService(repo repository, notifier notification.Service, publisher Publisher) Service {
	return &service{repo: repo, notifier: notifier, publisher: publisher}
}

func (s *service) ApproveEvent(ctx context.Context, eventID uuid.UUID) error {
	if err := s.repo.UpdateStatus(ctx, eventID, "approved"); err != nil {
		return err
	}

	if err := s.notifier.MarkAsReadByEventID(ctx, eventID); err != nil {
		log.Warn().Err(err).Msgf("Failed to mark notifications as read for event %s", eventID)
	}

	event, err := s.repo.LoadEvent(ctx, eventID)
	if err != nil {
		return err
	}

	_ = s.publisher.PublishStatusChanged(ctx, event)
	_ = s.publisher.Publish(ctx, event)

	return nil
}

func (s *service) RejectEvent(ctx context.Context, eventID uuid.UUID) error {
	if err := s.repo.UpdateStatus(ctx, eventID, "rejected"); err != nil {
		return err
	}

	// Mark related notifications as read
	if err := s.notifier.MarkAsReadByEventID(ctx, eventID); err != nil {
		log.Warn().Err(err).Msgf("Failed to mark notifications as read for event %s", eventID)
	}

	event, err := s.repo.LoadEvent(ctx, eventID)
	if err != nil {
		return err
	}

	_ = s.publisher.PublishStatusChanged(ctx, event)
	_ = s.publisher.Publish(ctx, event)

	return nil
}

func (s *service) GetEvent(ctx context.Context, eventID uuid.UUID) (events.Event, error) {
	return s.repo.LoadEvent(ctx, eventID)
}

func (s *service) GetEventTypes(ctx context.Context) ([]events.EventMeta, error) {
	metas := events.GetAllMetas()

	return metas, nil
}

func (s *service) LoadUserApprovedEvents(ctx context.Context, userID uuid.UUID) ([]events.Event, error) {
	return s.repo.LoadUserApprovedEvents(ctx, userID)
}

func (s *service) Load(ctx context.Context, id uuid.UUID) ([]events.Event, int, error) {
	return s.repo.Load(ctx, id)
}

func (s *service) Append(ctx context.Context, id uuid.UUID, expectedVersion int, newEvents ...events.Event) error {
	return s.repo.Append(ctx, id, expectedVersion, newEvents...)
}

func (s *service) UpdateStatus(ctx context.Context, eventID uuid.UUID, status string) error {
	return s.repo.UpdateStatus(ctx, eventID, status)
}
