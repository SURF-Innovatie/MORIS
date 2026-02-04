package events

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/notification"
	"github.com/SURF-Innovatie/MORIS/ent/user"
	orgsvc "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events/hydrator"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/projection"
	eventrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/event"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type ApprovalRequestNotificationHandler struct {
	cli       *ent.Client
	eventRepo *eventrepo.EntRepo
	rbac      orgsvc.Service
	hydrator  *hydrator.Hydrator
}

func NewApprovalRequestHandler(cli *ent.Client, eventRepo *eventrepo.EntRepo, rbac orgsvc.Service, hydrator *hydrator.Hydrator) *ApprovalRequestNotificationHandler {
	return &ApprovalRequestNotificationHandler{cli: cli, eventRepo: eventRepo, rbac: rbac, hydrator: hydrator}
}

func (h *ApprovalRequestNotificationHandler) Handle(ctx context.Context, e events.Event) error {
	if e.GetStatus() != "pending" {
		return nil
	}

	projectID := e.AggregateID()

	evts, _, err := h.eventRepo.Load(ctx, projectID)
	if err != nil || len(evts) == 0 {
		return err
	}

	proj := projection.Reduce(projectID, evts)
	if proj == nil || proj.OwningOrgNodeID == uuid.Nil {
		log.Warn().Msgf("ApprovalRequest: project %s has no owning org node", projectID)
		return nil
	}

	approvalNode, err := h.rbac.GetApprovalNode(ctx, proj.OwningOrgNodeID)
	if err != nil || approvalNode == nil {
		// no approver configured; you may want to log loudly
		log.Warn().Err(err).Msgf("ApprovalRequest: no approval node found for project %s", projectID)
		return nil
	}

	// Notify all admins effective for that approval node
	effs, err := h.rbac.ListEffectiveMemberships(ctx, approvalNode.ID)
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

		u, err := h.cli.User.Query().
			Where(user.PersonIDEQ(em.PersonID)).
			Only(ctx)
		if err != nil {
			continue
		}

		_, _ = h.cli.Notification.
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
	// Check if event implements Notifier
	if n, ok := e.(events.Notifier); ok {
		// Hydrate the event to get related entities for template variables
		de := h.hydrator.HydrateOne(ctx, e)
		// Add project title to the variables
		vars := n.NotificationVariables()
		if vars == nil {
			vars = make(map[string]string)
		}
		vars["project.Title"] = projectTitle
		vars = events.AddDetailedEventVariables(vars, de)
		msg := events.ResolveTemplate(n.ApprovalRequestTemplate(), vars)
		if msg != "" {
			return msg, nil
		}
	}

	return "", nil
}
