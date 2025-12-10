package event

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/projection"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	notifservice "github.com/SURF-Innovatie/MORIS/internal/notification"
)

type ProductAddedNotificationHandler struct {
	Notifier notifservice.Service
	Cli      *ent.Client
	ES       eventstore.Store
}

func (h *ProductAddedNotificationHandler) CanHandle(e events.Event) bool {
	_, ok := e.(events.ProductAdded)
	return ok
}

func (h *ProductAddedNotificationHandler) Handle(ctx context.Context, e events.Event) error {
	projectID := e.AggregateID()

	// Reconstruct project state to get current members
	evts, _, err := h.ES.Load(ctx, projectID)
	if err != nil {
		return err
	}

	if len(evts) == 0 {
		return nil
	}

	proj := projection.Reduce(projectID, evts)

	// Notify all current members
	for _, personID := range proj.People {
		// Find user for this person
		u, err := h.Cli.User.Query().
			Where(user.PersonIDEQ(personID)).
			Only(ctx)
		if err != nil {
			// Person might not have a user account or error
			continue
		}

		// Use the existing NotifyOfEvent logic (or replicate it here)
		// We can call Notifier.NotifyOfEvent if it existed, but we removed it.
		// So we implement it here.

		msg, err := h.buildMessage(ctx, e)
		if err != nil {
			continue
		}
		if msg == "" {
			continue
		}

		_, err = h.Cli.Notification.
			Create().
			SetMessage(msg).
			SetUser(u).
			SetEventID(e.GetID()).
			Save(ctx)
		if err != nil {
			// Log error?
		}
	}

	return nil
}

func (h *ProductAddedNotificationHandler) buildMessage(ctx context.Context, e events.Event) (string, error) {
	switch e.(type) {
	case events.ProductAdded:
		return "A new product has been added to the project.", nil
	default:
		return "", nil
	}
}
