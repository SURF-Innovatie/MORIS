package events

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/organisationnode"
	"github.com/SURF-Innovatie/MORIS/ent/user"
	events2 "github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events/hydrator"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/projection"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/event"
)

type ProjectEventNotificationHandler struct {
	cli       *ent.Client
	eventRepo *event.EntRepo
	hydrator  *hydrator.Hydrator
}

func NewProjectEventHandler(cli *ent.Client, eventRepo *event.EntRepo, hydrator *hydrator.Hydrator) *ProjectEventNotificationHandler {
	return &ProjectEventNotificationHandler{cli: cli, eventRepo: eventRepo, hydrator: hydrator}
}

func (h *ProjectEventNotificationHandler) Handle(ctx context.Context, e events2.Event) error {
	if e.GetStatus() == events2.StatusPending {
		return nil
	}

	msg, err := h.buildMessage(ctx, e)
	if err != nil || msg == "" {
		return err
	}

	// Reconstruct project state to get current members
	projectID := e.AggregateID()
	evts, _, err := h.eventRepo.Load(ctx, projectID)
	if err != nil {
		return err
	}

	if len(evts) == 0 {
		return nil
	}

	proj := projection.Reduce(projectID, evts)

	// Notify all current members
	for _, member := range proj.Members {
		// Find user for this person
		u, err := h.cli.User.Query().
			Where(user.PersonIDEQ(member.PersonID)).
			Only(ctx)
		if err != nil {
			// Person might not have a user account or error
			continue
		}

		_, err = h.cli.Notification.
			Create().
			SetMessage(msg).
			SetUser(u).
			SetEventID(e.GetID()).
			Save(ctx)
		if err != nil {
			continue
		}
	}

	return nil
}

func (h *ProjectEventNotificationHandler) buildMessage(ctx context.Context, e events2.Event) (string, error) {
	// Check if event implements Notifier
	if n, ok := e.(events2.Notifier); ok {
		// Hydrate the event to get related entities for template variables
		de := h.hydrator.HydrateOne(ctx, e)
		msg := events2.BuildNotificationMessage(n, de)
		if msg != "" {
			return msg, nil
		}
	}

	// Fallback for events requiring DB or special logic (e.g. OwningOrgNodeChanged)
	switch v := e.(type) {
	case *events2.OwningOrgNodeChanged:
		n, err := h.cli.OrganisationNode.
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
