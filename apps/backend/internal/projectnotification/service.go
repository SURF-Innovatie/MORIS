package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent/organisation"
	"github.com/SURF-Innovatie/MORIS/ent/person"
	"github.com/SURF-Innovatie/MORIS/ent/projectnotification"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type Service interface {
	Create(ctx context.Context, p entities.ProjectNotification) (*entities.ProjectNotification, error)
	Get(ctx context.Context, id uuid.UUID) (*entities.ProjectNotification, error)
	Update(ctx context.Context, id uuid.UUID, p entities.ProjectNotification) (*entities.ProjectNotification, error)
	List(ctx context.Context) ([]entities.ProjectNotification, error)
	NotifyForEvents(ctx context.Context, user *ent.User, projectID uuid.UUID, evts ...events.Event) error
}

type service struct {
	cli *ent.Client
}

func NewService(cli *ent.Client) Service {
	return &service{cli: cli}
}

func (s *service) Create(ctx context.Context, p entities.ProjectNotification) (*entities.ProjectNotification, error) {
	row, err := s.cli.ProjectNotification.
		Create().
		SetMessage(p.Message).
		SetUser(p.User).
		SetProjectID(p.ProjectId).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return mapRow(row, p.User), nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*entities.ProjectNotification, error) {
	row, err := s.cli.ProjectNotification.
		Query().
		Where(projectnotification.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	user, err := row.QueryUser().Only(ctx)
	if err != nil {
		return nil, err
	}
	return mapRow(row, user), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, p entities.ProjectNotification) (*entities.ProjectNotification, error) {
	row, err := s.cli.ProjectNotification.
		UpdateOneID(id).
		SetMessage(p.Message).
		SetUser(p.User).
		SetProjectID(p.ProjectId).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapRow(row, p.User), nil
}

func (s *service) List(ctx context.Context) ([]entities.ProjectNotification, error) {
	rows, err := s.cli.ProjectNotification.
		Query().
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]entities.ProjectNotification, 0, len(rows))
	for _, r := range rows {
		user, err := r.QueryUser().Only(ctx)
		if err != nil {
			return nil, err
		}
		out = append(out, *mapRow(r, user))
	}
	return out, nil
}

func mapRow(r *ent.ProjectNotification, u *ent.User) *entities.ProjectNotification {
	return &entities.ProjectNotification{
		Id:        r.ID,
		ProjectId: r.ProjectID,
		Message:   r.Message,
		User:      u,
	}
}

func (s *service) NotifyForEvents(
	ctx context.Context,
	user *ent.User,
	projectID uuid.UUID,
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
		_, err = s.cli.ProjectNotification.
			Create().
			SetID(uuid.New()). // optional if schema has Default(uuid.New)
			SetProjectID(projectID).
			SetMessage(msg).
			SetSentAt(time.Now().UTC()). // optional if schema already defaults
			SetUser(user).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	return nil
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
