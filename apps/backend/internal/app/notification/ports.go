package notification

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/notification"
	"github.com/google/uuid"
)

type NotificationRepository interface {
	Create(ctx context.Context, n notification.Notification) (*notification.Notification, error)
	Get(ctx context.Context, id uuid.UUID) (*notification.Notification, error)
	Update(ctx context.Context, id uuid.UUID, n notification.Notification) (*notification.Notification, error)
	List(ctx context.Context) ([]notification.Notification, error)
	ListForUser(ctx context.Context, userID uuid.UUID) ([]notification.Notification, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkAsReadByEventID(ctx context.Context, eventID uuid.UUID) error
}
