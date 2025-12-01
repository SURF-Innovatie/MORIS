package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent/notification"
	"github.com/SURF-Innovatie/MORIS/ent/organisation"
	"github.com/SURF-Innovatie/MORIS/ent/person"
	"github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"

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
	MarkAsRead(ctx context.Context, id uuid.UUID) error
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
		Read:    r.Read,
		User:    u,
		Event:   e,
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
			SetSentAt(time.Now().UTC()). // optional if schema already defaults
			SetUser(user).
			SetEventID(e.GetID()).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	return nil
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

	default:
		// unknown or uninteresting event type
		return "", nil
	}
}
