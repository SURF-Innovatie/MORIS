package notification

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent/event"
	"github.com/SURF-Innovatie/MORIS/ent/notification"
	"github.com/SURF-Innovatie/MORIS/ent/organisation"
	"github.com/SURF-Innovatie/MORIS/ent/person"
	"github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type Service interface {
	Create(ctx context.Context, p entities.Notification) (*entities.Notification, error)
	Get(ctx context.Context, id uuid.UUID) (*entities.Notification, error)
	Update(ctx context.Context, id uuid.UUID, p entities.Notification) (*entities.Notification, error)
	List(ctx context.Context) ([]entities.Notification, error)
	ListForUser(ctx context.Context, userID uuid.UUID) ([]entities.Notification, error)
	NotifyForEvents(ctx context.Context, user *ent.User, evts ...events.Event) error
	NotifyApprovers(ctx context.Context, evts ...events.Event) error
	NotifyStatusUpdate(ctx context.Context, event events.Event, status string) error
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	NotifyOfEvent(ctx context.Context, userId uuid.UUID, e events.Event) error
}

type service struct {
	cli *ent.Client
}

func NewService(cli *ent.Client) Service {
	return &service{cli: cli}
}

func (s *service) Create(ctx context.Context, p entities.Notification) (*entities.Notification, error) {
	row, err := s.cli.Notification.
		Create().
		SetMessage(p.Message).
		SetUser(p.User).
		SetEvent(p.Event).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return mapRow(row, p.User, p.Event), nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*entities.Notification, error) {
	row, err := s.cli.Notification.
		Query().
		Where(notification.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	user, err := row.QueryUser().Only(ctx)
	if err != nil {
		return nil, err
	}

	event, _ := row.QueryEvent().Only(ctx)

	return mapRow(row, user, event), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, p entities.Notification) (*entities.Notification, error) {
	row, err := s.cli.Notification.
		UpdateOneID(id).
		SetMessage(p.Message).
		SetUser(p.User).
		SetEvent(p.Event).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapRow(row, p.User, p.Event), nil
}

func (s *service) List(ctx context.Context) ([]entities.Notification, error) {
	rows, err := s.cli.Notification.
		Query().
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]entities.Notification, 0, len(rows))
	for _, r := range rows {
		user, err := r.QueryUser().Only(ctx)
		if err != nil {
			return nil, err
		}

		event, _ := r.QueryEvent().Only(ctx)

		out = append(out, *mapRow(r, user, event))
	}

	return out, nil
}

func (s *service) ListForUser(ctx context.Context, userID uuid.UUID) ([]entities.Notification, error) {
	rows, err := s.cli.Notification.
		Query().
		Where(notification.HasUserWith(user.IDEQ(userID))).
		Order(ent.Desc(notification.FieldSentAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]entities.Notification, 0, len(rows))
	for _, r := range rows {
		u, err := r.QueryUser().Only(ctx)
		if err != nil {
			return nil, err
		}

		// event can be empty
		event, _ := r.QueryEvent().Only(ctx)

		out = append(out, *mapRow(r, u, event))
	}
	return out, nil
}

func mapRow(r *ent.Notification, u *ent.User, e *ent.Event) *entities.Notification {
	return &entities.Notification{
		Id:      r.ID,
		Message: r.Message,
		Type:    r.Type.String(),
		Read:    r.Read,
		User:    u,
		Event:   e,
		SentAt:  r.SentAt,
	}
}

func (s *service) NotifyForEvents(
	ctx context.Context,
	user *ent.User,
	evts ...events.Event,
) error {
	if user == nil {
		return fmt.Errorf("user is nil in NotifyForEvents")
	}

	for _, e := range evts {
		msg, err := s.buildMessageFromEvent(ctx, e)
		if err != nil {
			return err
		}

		if msg == "" {
			// nothing to send for this event type
			continue
		}

		// Create a notification row
		_, err = s.cli.Notification.
			Create().
			SetID(uuid.New()). // optional if schema has Default(uuid.New)
			SetMessage(msg).
			SetUser(user).
			SetEventID(e.GetID()).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) NotifyApprovers(ctx context.Context, evts ...events.Event) error {
	for _, e := range evts {
		// Only notify for PersonAdded events for now
		if _, ok := e.(events.PersonAdded); !ok {
			continue
		}

		projectID := e.AggregateID()
		logrus.Infof("NotifyApprovers: Processing PersonAdded event for project %s", projectID)

		// Find the ProjectStarted event to get the admin
		// We query the event store (via ent client) for the ProjectStarted event of this project
		startedEvent, err := s.cli.Event.
			Query().
			Where(
				event.ProjectIDEQ(projectID),
				event.TypeEQ(events.ProjectStartedType),
			).
			First(ctx)
		if err != nil {
			logrus.Errorf("NotifyApprovers: Failed to find ProjectStarted event for project %s: %v", projectID, err)
			continue
		}

		// Get the payload to find the admin
		payload, err := startedEvent.QueryProjectStarted().Only(ctx)
		if err != nil {
			logrus.Errorf("NotifyApprovers: Failed to get ProjectStarted payload: %v", err)
			continue
		}

		adminID := payload.ProjectAdmin
		if adminID == uuid.Nil {
			logrus.Warnf("NotifyApprovers: Admin ID is nil for project %s", projectID)
			continue
		}

		// adminID is a PersonID, we need to find the User associated with this Person
		adminUser, err := s.cli.User.Query().
			Where(user.PersonIDEQ(adminID)).
			Only(ctx)
		if err != nil {
			logrus.Errorf("NotifyApprovers: Failed to find user for admin person %s: %v", adminID, err)
			continue
		}

		msg := fmt.Sprintf("Approval requested: Person added to project '%s'", payload.Title)
		logrus.Infof("NotifyApprovers: Sending notification to admin %s: %s", adminID, msg)

		_, err = s.cli.Notification.
			Create().
			SetMessage(msg).
			SetUser(adminUser).
			SetEventID(e.GetID()).
			SetType(notification.TypeApprovalRequest).
			Save(ctx)
		if err != nil {
		}
	}
	return nil
}

func (s *service) NotifyOfEvent(ctx context.Context, userId uuid.UUID, e events.Event) error {
	msg, err := s.buildMessageFromEvent(ctx, e)
	if err != nil {
		return err
	}

	if msg == "" {
		// nothing to send for this event type
		return nil
	}

	_, err = s.cli.Notification.
		Create().
		SetMessage(msg).
		SetUserID(userId).
		SetEventID(e.GetID()).
		SetType(notification.TypeInfo).
		Save(ctx)
	return err
}

func (s *service) NotifyStatusUpdate(ctx context.Context, event events.Event, status string) error {
	creatorID := event.CreatedByID()
	if creatorID == uuid.Nil {
		return nil
	}

	user, err := s.cli.User.Get(ctx, creatorID)
	if err != nil {
		// User might not exist anymore?
		return nil
	}

	msg := fmt.Sprintf("Your request '%s' has been %s.", event.Type(), status)

	_, err = s.cli.Notification.
		Create().
		SetMessage(msg).
		SetUser(user).
		SetEventID(event.GetID()).
		SetType(notification.TypeStatusUpdate).
		Save(ctx)

	return err
}

func (s *service) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	_, err := s.cli.Notification.UpdateOneID(id).SetRead(true).Save(ctx)
	return err
}

func (s *service) buildMessageFromEvent(ctx context.Context, e events.Event) (string, error) {
	switch v := e.(type) {

	case events.ProjectStarted:
		return fmt.Sprintf("Project '%s' has been started.", v.Title), nil

	case events.TitleChanged:
		return fmt.Sprintf("Project title changed to '%s'.", v.Title), nil

	case events.DescriptionChanged:
		return "Project description has been updated.", nil

	case events.StartDateChanged:
		return fmt.Sprintf("Project start date changed to %s.", v.StartDate.Format("2006-01-02")), nil

	case events.EndDateChanged:
		return fmt.Sprintf("Project end date changed to %s.", v.EndDate.Format("2006-01-02")), nil

	case events.OrganisationChanged:
		org, err := s.cli.Organisation.
			Query().
			Where(organisation.IDEQ(v.OrganisationID)).
			First(ctx)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Project organisation changed to '%s'.", org.Name), nil

	case events.PersonAdded:
		per, err := s.cli.Person.
			Query().
			Where(person.IDEQ(v.PersonId)).
			First(ctx)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Person %s was added to the project.", per.Name), nil

	case events.PersonRemoved:
		per, err := s.cli.Person.
			Query().
			Where(person.IDEQ(v.PersonId)).
			First(ctx)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Person %s was removed from the project.", per.Name), nil

	case events.ProductAdded:
		return "A new product has been added to the project.", nil

	case events.ProductRemoved:
		return "A product has been removed from the project.", nil
	default:
		// unknown or uninteresting event type
		return "", nil
	}
}
