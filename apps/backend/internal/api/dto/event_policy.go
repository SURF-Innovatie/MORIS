package dto

import "github.com/SURF-Innovatie/MORIS/internal/domain/entities"

// PolicyConditionDTO represents a condition for policy evaluation
type PolicyConditionDTO struct {
	Field    string `json:"field"`    // e.g. "event.title", "project.status"
	Operator string `json:"operator"` // "equals", "contains", "greater_than", etc.
	Value    any    `json:"value"`
}

// EventPolicyRequest is the request body for creating/updating an event policy
type EventPolicyRequest struct {
	Name                    string               `json:"name"`
	Description             *string              `json:"description,omitempty"`
	EventTypes              []string             `json:"event_types"`
	Conditions              []PolicyConditionDTO `json:"conditions,omitempty"`
	ActionType              string               `json:"action_type"` // "notify" | "request_approval"
	MessageTemplate         *string              `json:"message_template,omitempty"`
	RecipientUserIDs        []string             `json:"recipient_user_ids,omitempty"`
	RecipientProjectRoleIDs []string             `json:"recipient_project_role_ids,omitempty"`
	RecipientOrgRoleIDs     []string             `json:"recipient_org_role_ids,omitempty"`
	RecipientDynamic        []string             `json:"recipient_dynamic,omitempty"`
	Enabled                 bool                 `json:"enabled"`
}

// EventPolicyResponse is the response body for an event policy
type EventPolicyResponse struct {
	ID                      string               `json:"id"`
	Name                    string               `json:"name"`
	Description             *string              `json:"description,omitempty"`
	EventTypes              []string             `json:"event_types"`
	Conditions              []PolicyConditionDTO `json:"conditions,omitempty"`
	ActionType              string               `json:"action_type"`
	MessageTemplate         *string              `json:"message_template,omitempty"`
	RecipientUserIDs        []string             `json:"recipient_user_ids,omitempty"`
	RecipientProjectRoleIDs []string             `json:"recipient_project_role_ids,omitempty"`
	RecipientOrgRoleIDs     []string             `json:"recipient_org_role_ids,omitempty"`
	RecipientDynamic        []string             `json:"recipient_dynamic,omitempty"`
	OrgNodeID               *string              `json:"org_node_id,omitempty"`
	ProjectID               *string              `json:"project_id,omitempty"`
	Enabled                 bool                 `json:"enabled"`
	Inherited               bool                 `json:"inherited"`
	SourceOrgNodeID         *string              `json:"source_org_node_id,omitempty"`
	SourceOrgNodeName       *string              `json:"source_org_node_name,omitempty"`
}

// FromEntity converts domain entity to DTO
func (r *EventPolicyResponse) FromEntity(e *entities.EventPolicy) {
	if e == nil {
		return
	}

	r.ID = e.ID.String()
	r.Name = e.Name
	r.Description = e.Description
	r.EventTypes = e.EventTypes
	r.ActionType = string(e.ActionType)
	r.MessageTemplate = e.MessageTemplate
	r.RecipientDynamic = e.RecipientDynamic
	r.Enabled = e.Enabled
	r.Inherited = e.Inherited

	// Convert conditions
	r.Conditions = make([]PolicyConditionDTO, len(e.Conditions))
	for i, c := range e.Conditions {
		r.Conditions[i] = PolicyConditionDTO{
			Field:    c.Field,
			Operator: c.Operator,
			Value:    c.Value,
		}
	}

	// Convert UUIDs to strings
	for _, uid := range e.RecipientUserIDs {
		r.RecipientUserIDs = append(r.RecipientUserIDs, uid.String())
	}
	for _, rid := range e.RecipientProjectRoleIDs {
		r.RecipientProjectRoleIDs = append(r.RecipientProjectRoleIDs, rid.String())
	}
	for _, rid := range e.RecipientOrgRoleIDs {
		r.RecipientOrgRoleIDs = append(r.RecipientOrgRoleIDs, rid.String())
	}

	if e.OrgNodeID != nil {
		s := e.OrgNodeID.String()
		r.OrgNodeID = &s
	}
	if e.ProjectID != nil {
		s := e.ProjectID.String()
		r.ProjectID = &s
	}
	if e.SourceOrgNodeID != nil {
		s := e.SourceOrgNodeID.String()
		r.SourceOrgNodeID = &s
	}
	r.SourceOrgNodeName = e.SourceOrgNodeName
}
