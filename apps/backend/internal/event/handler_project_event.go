package event

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	notifservice "github.com/SURF-Innovatie/MORIS/internal/notification"
	"github.com/google/uuid"
)

type ProjectEventNotificationHandler struct {
	Notifier notifservice.Service
	Cli      *ent.Client
}

func (h *ProjectEventNotificationHandler) CanHandle(e events.Event) bool {
	switch e.(type) {
	case events.ProjectStarted,
		events.TitleChanged,
		events.DescriptionChanged,
		events.StartDateChanged,
		events.EndDateChanged,
		events.OrganisationChanged:
		return true
	}
	return false
}

func (h *ProjectEventNotificationHandler) Handle(ctx context.Context, e events.Event) error {
	creatorID := e.CreatedByID()
	if creatorID == uuid.Nil {
		return nil
	}

	user, err := h.Cli.User.Get(ctx, creatorID)
	if err != nil {
		return err
	}

	msg, err := h.buildMessage(ctx, e)
	if err != nil {
		return err
	}
	if msg == "" {
		return nil
	}

	_, err = h.Cli.Notification.
		Create().
		SetMessage(msg).
		SetUser(user).
		SetEventID(e.GetID()).
		Save(ctx)
	return err
}

func (h *ProjectEventNotificationHandler) buildMessage(ctx context.Context, e events.Event) (string, error) {
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
		org, err := h.Cli.Organisation.
			Query().
			Where(organisation.IDEQ(v.OrganisationID)).
			First(ctx)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Project organisation changed to '%s'.", org.Name), nil
	default:
		return "", nil
	}
}
