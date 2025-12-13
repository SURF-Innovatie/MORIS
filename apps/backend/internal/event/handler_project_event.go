package event

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/organisationnode"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
)

type ProjectEventNotificationHandler struct {
	Cli *ent.Client
}

func (h *ProjectEventNotificationHandler) CanHandle(e events.Event) bool {
	switch e.(type) {
	case events.ProjectStarted,
		events.TitleChanged,
		events.DescriptionChanged,
		events.StartDateChanged,
		events.EndDateChanged,
		events.OwningOrgNodeChanged:
		return true
	}
	return false
}

func (h *ProjectEventNotificationHandler) Handle(ctx context.Context, e events.Event) error {
	u, err := ResolveUser(ctx, h.Cli, e.CreatedByID())
	if err != nil || u == nil {
		return err
	}

	msg, err := h.buildMessage(ctx, e)
	if err != nil || msg == "" {
		return err
	}

	_, err = h.Cli.Notification.
		Create().
		SetMessage(msg).
		SetUser(u).
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
	case events.OwningOrgNodeChanged:
		n, err := h.Cli.OrganisationNode.
			Query().
			Where(organisationnode.IDEQ(v.OwningOrgNodeID)).
			Only(ctx)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Project owning organisation node changed to '%s'.", n.Name), nil
	default:
		return "", nil
	}
}
