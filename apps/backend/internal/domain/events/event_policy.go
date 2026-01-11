package events

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

const EventPolicyAddedType = "project.event_policy_added"
const EventPolicyRemovedType = "project.event_policy_removed"
const EventPolicyUpdatedType = "project.event_policy_updated"

// EventPolicyAdded represents adding an event policy to a project
type EventPolicyAdded struct {
	Base
	PolicyID                uuid.UUID   `json:"policy_id"`
	Name                    string      `json:"name"`
	Description             *string     `json:"description,omitempty"`
	EventTypes              []string    `json:"event_types"`
	ActionType              string      `json:"action_type"` // "notify" or "request_approval"
	RecipientUserIDs        []uuid.UUID `json:"recipient_user_ids,omitempty"`
	RecipientProjectRoleIDs []uuid.UUID `json:"recipient_project_role_ids,omitempty"`
	RecipientOrgRoleIDs     []uuid.UUID `json:"recipient_org_role_ids,omitempty"`
	RecipientDynamic        []string    `json:"recipient_dynamic,omitempty"`
	Enabled                 bool        `json:"enabled"`
}

func (EventPolicyAdded) isEvent()     {}
func (EventPolicyAdded) Type() string { return EventPolicyAddedType }
func (e EventPolicyAdded) String() string {
	return fmt.Sprintf("Event policy '%s' added", e.Name)
}

func (e *EventPolicyAdded) Apply(project *entities.Project) {
	// Policies are stored separately, not on project entity directly
}

func (e *EventPolicyAdded) NotificationMessage() string {
	return fmt.Sprintf("Event policy '%s' has been added to the project.", e.Name)
}

type EventPolicyAddedInput struct {
	Name                    string      `json:"name"`
	Description             *string     `json:"description,omitempty"`
	EventTypes              []string    `json:"event_types"`
	ActionType              string      `json:"action_type"`
	RecipientUserIDs        []uuid.UUID `json:"recipient_user_ids,omitempty"`
	RecipientProjectRoleIDs []uuid.UUID `json:"recipient_project_role_ids,omitempty"`
	RecipientOrgRoleIDs     []uuid.UUID `json:"recipient_org_role_ids,omitempty"`
	RecipientDynamic        []string    `json:"recipient_dynamic,omitempty"`
	Enabled                 bool        `json:"enabled"`
}

func DecideEventPolicyAdded(
	projectID uuid.UUID,
	actor uuid.UUID,
	in EventPolicyAddedInput,
	status Status,
) (Event, error) {
	return &EventPolicyAdded{
		Base:                    NewBase(projectID, actor, status),
		PolicyID:                uuid.New(), // Generate new policy ID
		Name:                    in.Name,
		Description:             in.Description,
		EventTypes:              in.EventTypes,
		ActionType:              in.ActionType,
		RecipientUserIDs:        in.RecipientUserIDs,
		RecipientProjectRoleIDs: in.RecipientProjectRoleIDs,
		RecipientOrgRoleIDs:     in.RecipientOrgRoleIDs,
		RecipientDynamic:        in.RecipientDynamic,
		Enabled:                 in.Enabled,
	}, nil
}

// EventPolicyRemoved represents removing an event policy from a project
type EventPolicyRemoved struct {
	Base
	PolicyID uuid.UUID `json:"policy_id"`
	Name     string    `json:"name"`
}

func (EventPolicyRemoved) isEvent()     {}
func (EventPolicyRemoved) Type() string { return EventPolicyRemovedType }
func (e EventPolicyRemoved) String() string {
	return fmt.Sprintf("Event policy '%s' removed", e.Name)
}

func (e *EventPolicyRemoved) Apply(project *entities.Project) {
	// Policies are stored separately
}

func (e *EventPolicyRemoved) NotificationMessage() string {
	return fmt.Sprintf("Event policy '%s' has been removed from the project.", e.Name)
}

type EventPolicyRemovedInput struct {
	PolicyID uuid.UUID `json:"policy_id"`
	Name     string    `json:"name"`
}

func DecideEventPolicyRemoved(
	projectID uuid.UUID,
	actor uuid.UUID,
	in EventPolicyRemovedInput,
	status Status,
) (Event, error) {
	return &EventPolicyRemoved{
		Base:     NewBase(projectID, actor, status),
		PolicyID: in.PolicyID,
		Name:     in.Name,
	}, nil
}

// EventPolicyUpdated represents modifying an existing event policy on a project
type EventPolicyUpdated struct {
	Base
	PolicyID                uuid.UUID   `json:"policy_id"`
	Name                    string      `json:"name"`
	Description             *string     `json:"description,omitempty"`
	EventTypes              []string    `json:"event_types"`
	ActionType              string      `json:"action_type"`
	RecipientUserIDs        []uuid.UUID `json:"recipient_user_ids,omitempty"`
	RecipientProjectRoleIDs []uuid.UUID `json:"recipient_project_role_ids,omitempty"`
	RecipientOrgRoleIDs     []uuid.UUID `json:"recipient_org_role_ids,omitempty"`
	RecipientDynamic        []string    `json:"recipient_dynamic,omitempty"`
	Enabled                 bool        `json:"enabled"`
}

func (EventPolicyUpdated) isEvent()     {}
func (EventPolicyUpdated) Type() string { return EventPolicyUpdatedType }
func (e EventPolicyUpdated) String() string {
	return fmt.Sprintf("Event policy '%s' updated", e.Name)
}

func (e *EventPolicyUpdated) Apply(project *entities.Project) {
	// Policies are stored separately
}

func (e *EventPolicyUpdated) NotificationMessage() string {
	return fmt.Sprintf("Event policy '%s' has been updated.", e.Name)
}

type EventPolicyUpdatedInput struct {
	PolicyID                uuid.UUID   `json:"policy_id"`
	Name                    string      `json:"name"`
	Description             *string     `json:"description,omitempty"`
	EventTypes              []string    `json:"event_types"`
	ActionType              string      `json:"action_type"`
	RecipientUserIDs        []uuid.UUID `json:"recipient_user_ids,omitempty"`
	RecipientProjectRoleIDs []uuid.UUID `json:"recipient_project_role_ids,omitempty"`
	RecipientOrgRoleIDs     []uuid.UUID `json:"recipient_org_role_ids,omitempty"`
	RecipientDynamic        []string    `json:"recipient_dynamic,omitempty"`
	Enabled                 bool        `json:"enabled"`
}

func DecideEventPolicyUpdated(
	projectID uuid.UUID,
	actor uuid.UUID,
	in EventPolicyUpdatedInput,
	status Status,
) (Event, error) {
	return &EventPolicyUpdated{
		Base:                    NewBase(projectID, actor, status),
		PolicyID:                in.PolicyID,
		Name:                    in.Name,
		Description:             in.Description,
		EventTypes:              in.EventTypes,
		ActionType:              in.ActionType,
		RecipientUserIDs:        in.RecipientUserIDs,
		RecipientProjectRoleIDs: in.RecipientProjectRoleIDs,
		RecipientOrgRoleIDs:     in.RecipientOrgRoleIDs,
		RecipientDynamic:        in.RecipientDynamic,
		Enabled:                 in.Enabled,
	}, nil
}

func init() {
	// Register EventPolicyAdded
	RegisterMeta(EventMeta{
		Type:         EventPolicyAddedType,
		FriendlyName: "Event Policy Added",
	}, func() Event { return &EventPolicyAdded{} })

	RegisterDecider[EventPolicyAddedInput](EventPolicyAddedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *entities.Project, in EventPolicyAddedInput, status Status) (Event, error) {
			return DecideEventPolicyAdded(projectID, actor, in, status)
		})

	RegisterInputType(EventPolicyAddedType, EventPolicyAddedInput{})

	// Register EventPolicyRemoved
	RegisterMeta(EventMeta{
		Type:         EventPolicyRemovedType,
		FriendlyName: "Event Policy Removed",
	}, func() Event { return &EventPolicyRemoved{} })

	RegisterDecider[EventPolicyRemovedInput](EventPolicyRemovedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *entities.Project, in EventPolicyRemovedInput, status Status) (Event, error) {
			return DecideEventPolicyRemoved(projectID, actor, in, status)
		})

	RegisterInputType(EventPolicyRemovedType, EventPolicyRemovedInput{})

	// Register EventPolicyUpdated
	RegisterMeta(EventMeta{
		Type:         EventPolicyUpdatedType,
		FriendlyName: "Event Policy Updated",
	}, func() Event { return &EventPolicyUpdated{} })

	RegisterDecider[EventPolicyUpdatedInput](EventPolicyUpdatedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *entities.Project, in EventPolicyUpdatedInput, status Status) (Event, error) {
			return DecideEventPolicyUpdated(projectID, actor, in, status)
		})

	RegisterInputType(EventPolicyUpdatedType, EventPolicyUpdatedInput{})
}
