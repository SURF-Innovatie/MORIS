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
	cli *ent.Client
	es  eventstore.Store
}

func NewProjectEventHandler(cli *ent.Client, es eventstore.Store) *ProjectEventNotificationHandler {
	return &ProjectEventNotificationHandler{cli: cli, es: es}
}

func (h *ProjectEventNotificationHandler) Handle(ctx context.Context, e events.Event) error {
	if e.GetStatus() == events.StatusPending {
		return nil
	}
	// First check metadata policy
	meta := events.GetMeta(e.Type())
	if !meta.ShouldNotify(ctx, e, h.cli) {
		return nil
	}

	msg, err := h.buildMessage(ctx, e)
	if err != nil || msg == "" {
		return err
	}

	// Reconstruct project state to get current members
	projectID := e.AggregateID()
	evts, _, err := h.es.Load(ctx, projectID)
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

func (h *ProjectEventNotificationHandler) buildMessage(ctx context.Context, e events.Event) (string, error) {
	// Check if event implements Notifier
	if n, ok := e.(events.Notifier); ok {
		return n.NotificationMessage(), nil
	}

	// Fallback for events requiring DB or special logic (e.g. OwningOrgNodeChanged)
	switch v := e.(type) {
	case *events.OwningOrgNodeChanged:
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
