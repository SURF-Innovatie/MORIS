package event

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/app/eventpolicy"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/sirupsen/logrus"
)

// EventPolicyHandler handles EventPolicyAdded, EventPolicyRemoved, and EventPolicyUpdated events
// to persist/update/remove policies when these events are processed
type EventPolicyHandler struct {
	PolicyRepo eventpolicy.Repository
	Cli        *ent.Client
}

func (h *EventPolicyHandler) Handle(ctx context.Context, event events.Event) error {
	switch e := event.(type) {
	case *events.EventPolicyAdded:
		return h.handlePolicyAdded(ctx, e)
	case *events.EventPolicyRemoved:
		return h.handlePolicyRemoved(ctx, e)
	case *events.EventPolicyUpdated:
		return h.handlePolicyUpdated(ctx, e)
	}
	return nil
}

func (h *EventPolicyHandler) handlePolicyAdded(ctx context.Context, e *events.EventPolicyAdded) error {
	projectID := e.AggregateID()

	policy := entities.EventPolicy{
		ID:                      e.PolicyID,
		Name:                    e.Name,
		Description:             e.Description,
		EventTypes:              e.EventTypes,
		ActionType:              entities.ActionType(e.ActionType),
		RecipientUserIDs:        e.RecipientUserIDs,
		RecipientProjectRoleIDs: e.RecipientProjectRoleIDs,
		RecipientOrgRoleIDs:     e.RecipientOrgRoleIDs,
		RecipientDynamic:        e.RecipientDynamic,
		ProjectID:               &projectID,
		Enabled:                 e.Enabled,
	}

	_, err := h.PolicyRepo.Create(ctx, policy)
	if err != nil {
		logrus.Errorf("Failed to create policy from event: %v", err)
		return err
	}

	logrus.Infof("Event policy '%s' created for project %s", e.Name, projectID)
	return nil
}

func (h *EventPolicyHandler) handlePolicyRemoved(ctx context.Context, e *events.EventPolicyRemoved) error {
	err := h.PolicyRepo.Delete(ctx, e.PolicyID)
	if err != nil {
		logrus.Errorf("Failed to delete policy from event: %v", err)
		return err
	}

	logrus.Infof("Event policy '%s' removed from project %s", e.Name, e.AggregateID())
	return nil
}

func (h *EventPolicyHandler) handlePolicyUpdated(ctx context.Context, e *events.EventPolicyUpdated) error {
	projectID := e.AggregateID()

	policy := entities.EventPolicy{
		ID:                      e.PolicyID,
		Name:                    e.Name,
		Description:             e.Description,
		EventTypes:              e.EventTypes,
		ActionType:              entities.ActionType(e.ActionType),
		RecipientUserIDs:        e.RecipientUserIDs,
		RecipientProjectRoleIDs: e.RecipientProjectRoleIDs,
		RecipientOrgRoleIDs:     e.RecipientOrgRoleIDs,
		RecipientDynamic:        e.RecipientDynamic,
		ProjectID:               &projectID,
		Enabled:                 e.Enabled,
	}

	_, err := h.PolicyRepo.Update(ctx, e.PolicyID, policy)
	if err != nil {
		logrus.Errorf("Failed to update policy from event: %v", err)
		return err
	}

	logrus.Infof("Event policy '%s' updated for project %s", e.Name, projectID)
	return nil
}
