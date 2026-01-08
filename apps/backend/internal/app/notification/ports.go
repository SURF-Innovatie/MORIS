package notification

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type NotificationRepository interface {
	Create(ctx context.Context, n entities.Notification) (*entities.Notification, error)
	Get(ctx context.Context, id uuid.UUID) (*entities.Notification, error)
	Update(ctx context.Context, id uuid.UUID, n entities.Notification) (*entities.Notification, error)
	List(ctx context.Context) ([]entities.Notification, error)
	ListForUser(ctx context.Context, userID uuid.UUID) ([]entities.Notification, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
}
