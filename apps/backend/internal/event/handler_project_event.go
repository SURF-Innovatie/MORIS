package event

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/organisationnode"
	"github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/projection"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
)

type ProjectEventNotificationHandler struct {
	Cli *ent.Client
	ES  eventstore.Store
}

func (h *ProjectEventNotificationHandler) Handle(ctx context.Context, e events.Event) error {
	// First check metadata policy
	meta := events.GetMeta(e.Type())
	if !meta.ShouldNotify(ctx, e, h.Cli) {
		return nil
	}

	msg, err := h.buildMessage(ctx, e)
	if err != nil || msg == "" {
		return err
	}

	// Reconstruct project state to get current members
	projectID := e.AggregateID()
	evts, _, err := h.ES.Load(ctx, projectID)
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
		u, err := h.Cli.User.Query().
			Where(user.PersonIDEQ(member.PersonID)).
			Only(ctx)
		if err != nil {
			// Person might not have a user account or error
			continue
		}

		_, err = h.Cli.Notification.
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

func (h *ProjectEventNotificationHandler) buildMessage(ctx context.Context, e events.Event) (string, error) {
	// Check if event implements Notifier
	if n, ok := e.(events.Notifier); ok {
		return n.NotificationMessage(), nil
	}

	// Fallback for events requiring DB or special logic (e.g. OwningOrgNodeChanged)
	switch v := e.(type) {
	case *events.OwningOrgNodeChanged:
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
