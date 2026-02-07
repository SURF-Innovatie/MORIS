package notification

import (
	"context"
	"encoding/json"
	"os"

	"github.com/SURF-Innovatie/MORIS/internal/domain/ldn"
	"github.com/SURF-Innovatie/MORIS/internal/domain/notification"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, n notification.Notification) (*notification.Notification, error)
	Send(ctx context.Context, userIDs []uuid.UUID, eventID uuid.UUID, message string, notificationType notification.NotificationType) error
	Get(ctx context.Context, id uuid.UUID) (*notification.Notification, error)
	Update(ctx context.Context, id uuid.UUID, n notification.Notification) (*notification.Notification, error)
	List(ctx context.Context) ([]notification.Notification, error)
	ListForUser(ctx context.Context, userID uuid.UUID) ([]notification.Notification, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkAsReadByEventID(ctx context.Context, eventID uuid.UUID) error

	// LDN/AS2 methods
	CreateFromActivity(ctx context.Context, activity *ldn.Activity, userID uuid.UUID) (*notification.Notification, error)
	GetAsActivity(ctx context.Context, id uuid.UUID) (*ldn.Activity, error)
	ListForUserAsActivities(ctx context.Context, userID uuid.UUID) ([]ldn.Activity, error)
}

type service struct {
	repo NotificationRepository
}

func NewService(repo NotificationRepository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, n notification.Notification) (*notification.Notification, error) {
	return s.repo.Create(ctx, n)
}

func (s *service) Send(ctx context.Context, userIDs []uuid.UUID, eventID uuid.UUID, message string, notificationType notification.NotificationType) error {
	if len(userIDs) == 0 {
		return nil
	}

	// Generate AS2 activity for the notification
	originURL := os.Getenv("LDN_ORIGIN_URL")
	if originURL == "" {
		originURL = "http://localhost:8080"
	}

	as2Type := ldn.AS2TypeFromLegacy(string(notificationType))

	for _, userID := range userIDs {
		// Build AS2 activity
		activity := ldn.NewActivity(as2Type, originURL)
		activity.Summary = message
		activity.Object = ldn.NewObject("urn:uuid:"+eventID.String(), "Event")

		payloadBytes, _ := json.Marshal(activity)
		payloadStr := string(payloadBytes)
		activityID := activity.ID

		n := notification.Notification{
			UserID:         userID,
			Message:        message,
			EventID:        &eventID,
			Type:           notificationType,
			ActivityID:     &activityID,
			ActivityType:   &as2Type,
			Payload:        &payloadStr,
			OriginService:  &originURL,
			Direction:      notification.DirectionInternal,
			DeliveryStatus: notification.DeliveryDelivered,
		}

		if _, err := s.repo.Create(ctx, n); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*notification.Notification, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, n notification.Notification) (*notification.Notification, error) {
	return s.repo.Update(ctx, id, n)
}

func (s *service) List(ctx context.Context) ([]notification.Notification, error) {
	return s.repo.List(ctx)
}

func (s *service) ListForUser(ctx context.Context, userID uuid.UUID) ([]notification.Notification, error) {
	return s.repo.ListForUser(ctx, userID)
}

func (s *service) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	return s.repo.MarkAsRead(ctx, id)
}

func (s *service) MarkAsReadByEventID(ctx context.Context, eventID uuid.UUID) error {
	return s.repo.MarkAsReadByEventID(ctx, eventID)
}

// CreateFromActivity creates a notification from an incoming LDN Activity.
func (s *service) CreateFromActivity(ctx context.Context, activity *ldn.Activity, userID uuid.UUID) (*notification.Notification, error) {
	payloadBytes, err := json.Marshal(activity)
	if err != nil {
		return nil, err
	}
	payloadStr := string(payloadBytes)

	// Map AS2 type to legacy type
	legacyType := notification.NotificationType(ldn.LegacyTypeFromAS2(activity.Type[0]))

	var originService, targetService *string
	if activity.Origin != nil {
		originService = &activity.Origin.ID
	}
	if activity.Target != nil {
		targetService = &activity.Target.ID
	}

	activityType := ""
	if len(activity.Type) > 0 {
		activityType = activity.Type[0]
	}

	n := notification.Notification{
		UserID:         userID,
		Message:        activity.Summary,
		Type:           legacyType,
		ActivityID:     &activity.ID,
		ActivityType:   &activityType,
		Payload:        &payloadStr,
		OriginService:  originService,
		TargetService:  targetService,
		Direction:      notification.DirectionInbound,
		DeliveryStatus: notification.DeliveryDelivered,
	}

	return s.repo.Create(ctx, n)
}

// GetAsActivity returns a notification as an AS2 Activity.
func (s *service) GetAsActivity(ctx context.Context, id uuid.UUID) (*ldn.Activity, error) {
	n, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return notificationToActivity(n)
}

// ListForUserAsActivities returns notifications as AS2 Activities.
func (s *service) ListForUserAsActivities(ctx context.Context, userID uuid.UUID) ([]ldn.Activity, error) {
	notifications, err := s.repo.ListForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	activities := make([]ldn.Activity, 0, len(notifications))
	for _, n := range notifications {
		activity, err := notificationToActivity(&n)
		if err != nil {
			continue // Skip notifications that can't be converted
		}
		activities = append(activities, *activity)
	}
	return activities, nil
}

// notificationToActivity converts a notification to an AS2 Activity.
func notificationToActivity(n *notification.Notification) (*ldn.Activity, error) {
	// If we have a stored payload, deserialize it
	if n.Payload != nil && *n.Payload != "" {
		var activity ldn.Activity
		if err := json.Unmarshal([]byte(*n.Payload), &activity); err == nil {
			return &activity, nil
		}
	}

	// Otherwise, construct from fields
	originURL := os.Getenv("LDN_ORIGIN_URL")
	if originURL == "" {
		originURL = "http://localhost:8080"
	}

	as2Type := ldn.AS2TypeFromLegacy(string(n.Type))
	activity := ldn.NewActivity(as2Type, originURL)
	activity.ID = "urn:uuid:" + n.ID.String()
	activity.Summary = n.Message
	activity.Published = &n.SentAt

	if n.EventID != nil {
		activity.Object = ldn.NewObject("urn:uuid:"+n.EventID.String(), "Event")
	}

	return activity, nil
}
