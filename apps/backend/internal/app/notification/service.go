package notification

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, n entities.Notification) (*entities.Notification, error)
	Send(ctx context.Context, userIDs []uuid.UUID, eventID uuid.UUID, message string, notificationType entities.NotificationType) error
	Get(ctx context.Context, id uuid.UUID) (*entities.Notification, error)
	Update(ctx context.Context, id uuid.UUID, n entities.Notification) (*entities.Notification, error)
	List(ctx context.Context) ([]entities.Notification, error)
	ListForUser(ctx context.Context, userID uuid.UUID) ([]entities.Notification, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkAsReadByEventID(ctx context.Context, eventID uuid.UUID) error
}

type service struct {
	repo NotificationRepository
}

func NewService(repo NotificationRepository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, n entities.Notification) (*entities.Notification, error) {
	return s.repo.Create(ctx, n)
}

func (s *service) Send(ctx context.Context, userIDs []uuid.UUID, eventID uuid.UUID, message string, notificationType entities.NotificationType) error {
	if len(userIDs) == 0 {
		return nil
	}

	for _, userID := range userIDs {
		n := entities.Notification{
			UserID:  userID,
			Message: message,
			EventID: &eventID,
			Type:    notificationType,
		}

		if _, err := s.repo.Create(ctx, n); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*entities.Notification, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, n entities.Notification) (*entities.Notification, error) {
	return s.repo.Update(ctx, id, n)
}

func (s *service) List(ctx context.Context) ([]entities.Notification, error) {
	return s.repo.List(ctx)
}

func (s *service) ListForUser(ctx context.Context, userID uuid.UUID) ([]entities.Notification, error) {
	return s.repo.ListForUser(ctx, userID)
}

func (s *service) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	return s.repo.MarkAsRead(ctx, id)
}

func (s *service) MarkAsReadByEventID(ctx context.Context, eventID uuid.UUID) error {
	return s.repo.MarkAsReadByEventID(ctx, eventID)
}
