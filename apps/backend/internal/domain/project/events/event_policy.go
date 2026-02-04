package events

import (
	"context"
	"fmt"

	projdomain "github.com/SURF-Innovatie/MORIS/internal/domain/project"
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

func (e *EventPolicyAdded) Apply(p *projdomain.Project) {
	// Policies are stored separately, not on project entity directly
}

func (e *EventPolicyAdded) NotificationTemplate() string {
	return "Event policy '{{event.Name}}' has been added to the project."
}

func (e *EventPolicyAdded) ApprovalRequestTemplate() string {
	return "Adding event policy '{{event.Name}}' requires approval."
}

func (e *EventPolicyAdded) ApprovedTemplate() string {
	return "Event policy '{{event.Name}}' has been approved and added."
}

func (e *EventPolicyAdded) RejectedTemplate() string {
	return "Adding event policy '{{event.Name}}' has been rejected."
}

func (e *EventPolicyAdded) NotificationVariables() map[string]string {
	return map[string]string{
		"event.Name": e.Name,
	}
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
	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = EventPolicyAddedMeta.FriendlyName

	return &EventPolicyAdded{
		Base:                    base,
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

func (e *EventPolicyRemoved) Apply(p *projdomain.Project) {
	// Policies are stored separately
}

func (e *EventPolicyRemoved) NotificationTemplate() string {
	return "Event policy '{{event.Name}}' has been removed from the project."
}

func (e *EventPolicyRemoved) ApprovalRequestTemplate() string {
	return "Removing event policy '{{event.Name}}' requires approval."
}

func (e *EventPolicyRemoved) ApprovedTemplate() string {
	return "Event policy '{{event.Name}}' removal has been approved."
}

func (e *EventPolicyRemoved) RejectedTemplate() string {
	return "Removing event policy '{{event.Name}}' has been rejected."
}

func (e *EventPolicyRemoved) NotificationVariables() map[string]string {
	return map[string]string{
		"event.Name": e.Name,
	}
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
	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = EventPolicyRemovedMeta.FriendlyName

	return &EventPolicyRemoved{
		Base:     base,
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

func (e *EventPolicyUpdated) Apply(p *projdomain.Project) {
	// Policies are stored separately
}

func (e *EventPolicyUpdated) NotificationTemplate() string {
	return "Event policy '{{event.Name}}' has been updated."
}

func (e *EventPolicyUpdated) ApprovalRequestTemplate() string {
	return "Updating event policy '{{event.Name}}' requires approval."
}

func (e *EventPolicyUpdated) ApprovedTemplate() string {
	return "Event policy '{{event.Name}}' update has been approved."
}

func (e *EventPolicyUpdated) RejectedTemplate() string {
	return "Updating event policy '{{event.Name}}' has been rejected."
}

func (e *EventPolicyUpdated) NotificationVariables() map[string]string {
	return map[string]string{
		"event.Name": e.Name,
	}
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
	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = EventPolicyUpdatedMeta.FriendlyName

	return &EventPolicyUpdated{
		Base:                    base,
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

var EventPolicyAddedMeta = EventMeta{
	Type:         EventPolicyAddedType,
	FriendlyName: "Event Policy Added",
}
var EventPolicyRemovedMeta = EventMeta{
	Type:         EventPolicyRemovedType,
	FriendlyName: "Event Policy Removed",
}
var EventPolicyUpdatedMeta = EventMeta{
	Type:         EventPolicyUpdatedType,
	FriendlyName: "Event Policy Updated",
}

func init() {
	// Register EventPolicyAdded
	RegisterMeta(EventPolicyAddedMeta, func() Event {
		return &EventPolicyAdded{
			Base: Base{FriendlyNameStr: EventPolicyAddedMeta.FriendlyName},
		}
	})

	RegisterDecider[EventPolicyAddedInput](EventPolicyAddedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *projdomain.Project, in EventPolicyAddedInput, status Status) (Event, error) {
			return DecideEventPolicyAdded(projectID, actor, in, status)
		})

	RegisterInputType(EventPolicyAddedType, EventPolicyAddedInput{})

	// Register EventPolicyRemoved
	RegisterMeta(EventPolicyRemovedMeta, func() Event {
		return &EventPolicyRemoved{
			Base: Base{FriendlyNameStr: EventPolicyRemovedMeta.FriendlyName},
		}
	})

	RegisterDecider[EventPolicyRemovedInput](EventPolicyRemovedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *projdomain.Project, in EventPolicyRemovedInput, status Status) (Event, error) {
			return DecideEventPolicyRemoved(projectID, actor, in, status)
		})

	RegisterInputType(EventPolicyRemovedType, EventPolicyRemovedInput{})

	// Register EventPolicyUpdated
	RegisterMeta(EventPolicyUpdatedMeta, func() Event {
		return &EventPolicyUpdated{
			Base: Base{FriendlyNameStr: EventPolicyUpdatedMeta.FriendlyName},
		}
	})

	RegisterDecider[EventPolicyUpdatedInput](EventPolicyUpdatedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *projdomain.Project, in EventPolicyUpdatedInput, status Status) (Event, error) {
			return DecideEventPolicyUpdated(projectID, actor, in, status)
		})

	RegisterInputType(EventPolicyUpdatedType, EventPolicyUpdatedInput{})
}
