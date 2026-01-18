package event

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/notification"
	"github.com/SURF-Innovatie/MORIS/ent/user"
	orgsvc "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/projection"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ApprovalRequestNotificationHandler struct {
	Cli  *ent.Client
	ES   eventstore.Store
	RBAC orgsvc.Service
}

func (h *ApprovalRequestNotificationHandler) Handle(ctx context.Context, e events.Event) error {
	// Legacy approval logic removed. Policies now handle this.
	return nil

	projectID := e.AggregateID()

	evts, _, err := h.ES.Load(ctx, projectID)
	if err != nil || len(evts) == 0 {
		return err
	}

	proj := projection.Reduce(projectID, evts)
	if proj == nil || proj.OwningOrgNodeID == uuid.Nil {
		logrus.Warnf("ApprovalRequest: project %s has no owning org node", projectID)
		return nil
	}

	approvalNode, err := h.RBAC.GetApprovalNode(ctx, proj.OwningOrgNodeID)
	if err != nil || approvalNode == nil {
		// no approver configured; you may want to log loudly
		logrus.Warnf("ApprovalRequest: no approval node found for project %s: %v", projectID, err)
		return nil
	}

	// Notify all admins effective for that approval node
	effs, err := h.RBAC.ListEffectiveMemberships(ctx, approvalNode.ID)
	if err != nil {
		return err
	}

	msg, err := h.buildApprovalMessage(ctx, e, proj.Title)
	if err != nil {
		return err
	}
	if msg == "" {
		return nil
	}

	for _, em := range effs {
		hasAdminRights := false
		for _, p := range em.Permissions {
			if p == "manage_details" {
				hasAdminRights = true
				break
			}
		}
		if !hasAdminRights {
			continue
		}

		u, err := h.Cli.User.Query().
			Where(user.PersonIDEQ(em.PersonID)).
			Only(ctx)
		if err != nil {
			continue
		}

		_, _ = h.Cli.Notification.
			Create().
			SetMessage(msg).
			SetUser(u).
			SetEventID(e.GetID()).
			SetType(notification.TypeApprovalRequest).
			Save(ctx)
	}

	return nil
}

func (h *ApprovalRequestNotificationHandler) buildApprovalMessage(ctx context.Context, e events.Event, projectTitle string) (string, error) {
	if n, ok := e.(events.ApprovalNotifier); ok {
		return n.ApprovalMessage(projectTitle), nil
	}
	return "", nil
}
