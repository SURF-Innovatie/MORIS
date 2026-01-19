package event

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/app/eventpolicy"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/rs/zerolog/log"
)

// EventPolicyHandler handles EventPolicyAdded, EventPolicyRemoved, and EventPolicyUpdated events
// to persist/update/remove policies when these events are processed
type EventPolicyHandler struct {
	policyRepo eventpolicy.Repository
	cli        *ent.Client
}

func NewEventPolicyHandler(policyRepo eventpolicy.Repository, cli *ent.Client) *EventPolicyHandler {
	return &EventPolicyHandler{policyRepo: policyRepo, cli: cli}
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

	_, err := h.policyRepo.Create(ctx, policy)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create policy from event")
		return err
	}

	log.Info().Msgf("Event policy '%s' created for project %s", e.Name, projectID)
	return nil
}

func (h *EventPolicyHandler) handlePolicyRemoved(ctx context.Context, e *events.EventPolicyRemoved) error {
	err := h.policyRepo.Delete(ctx, e.PolicyID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete policy from event")
		return err
	}

	log.Info().Msgf("Event policy '%s' removed from project %s", e.Name, e.AggregateID())
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

	_, err := h.policyRepo.Update(ctx, e.PolicyID, policy)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update policy from event")
		return err
	}

	log.Info().Msgf("Event policy '%s' updated for project %s", e.Name, projectID)
	return nil
}
