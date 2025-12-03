package event

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	"github.com/SURF-Innovatie/MORIS/internal/notification"
	"github.com/google/uuid"
)

type Service interface {
	ApproveEvent(ctx context.Context, eventID uuid.UUID) error
	RejectEvent(ctx context.Context, eventID uuid.UUID) error
}

type service struct {
	cli      *ent.Client
	es       eventstore.Store
	notifier notification.Service
}

func NewService(es eventstore.Store, cli *ent.Client, notifier notification.Service) Service {
	return &service{es: es, cli: cli, notifier: notifier}
}

func (s *service) ApproveEvent(ctx context.Context, eventID uuid.UUID) error {
	// 1. Update status
	if err := s.es.UpdateEventStatus(ctx, eventID, "approved"); err != nil {
		return err
	}

	// 2. Notify creator
	event, err := s.es.LoadEvent(ctx, eventID)
	if err != nil {
		// Log error but don't fail the request?
		// Or fail?
		return err
	}

	return s.notifier.NotifyStatusUpdate(ctx, event, "approved")
}

func (s *service) RejectEvent(ctx context.Context, eventID uuid.UUID) error {
	if err := s.es.UpdateEventStatus(ctx, eventID, "rejected"); err != nil {
		return err
	}

	event, err := s.es.LoadEvent(ctx, eventID)
	if err != nil {
		return err
	}

	return s.notifier.NotifyStatusUpdate(ctx, event, "rejected")
}
