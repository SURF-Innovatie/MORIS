package policy

import (
	"slices"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

// ActionType represents the type of action a policy performs
type ActionType string

const (
	ActionTypeNotify          ActionType = "notify"
	ActionTypeRequestApproval ActionType = "request_approval"
)

// Condition represents a single condition for policy evaluation.
// Conditions are AND-ed together; all must pass for the policy to trigger.
type Condition struct {
	Field    string `json:"field"`    // Path to field: "event.<field>", "project.<field>", "custom_field.<name>"
	Operator string `json:"operator"` // Operator type (see constants below)
	Value    any    `json:"value"`    // Comparison value
}

// Supported condition operators (extensible)
const (
	OperatorEquals      = "equals"
	OperatorNotEquals   = "not_equals"
	OperatorContains    = "contains"
	OperatorStartsWith  = "starts_with"
	OperatorGreaterThan = "greater_than"
	OperatorLessThan    = "less_than"
	OperatorBetween     = "between"
	OperatorIn          = "in"
	OperatorNotIn       = "not_in"
	OperatorExists      = "exists"
	OperatorNotExists   = "not_exists"
)

// EventPolicy represents a configurable event trigger with actions.
type EventPolicy struct {
	ID                      uuid.UUID
	Name                    string
	Description             *string
	EventTypes              []string
	Conditions              []Condition
	ActionType              ActionType
	MessageTemplate         *string
	RecipientUserIDs        []uuid.UUID
	RecipientProjectRoleIDs []uuid.UUID
	RecipientOrgRoleIDs     []uuid.UUID
	RecipientDynamic        []string // "project_members", "project_owner", "org_admins"
	OrgNodeID               *uuid.UUID
	ProjectID               *uuid.UUID
	Enabled                 bool

	// Inheritance info (populated when querying with inheritance context)
	Inherited         bool
	SourceOrgNodeID   *uuid.UUID
	SourceOrgNodeName *string
}

// FromEnt converts an ent EventPolicy to domain entity
func (e *EventPolicy) FromEnt(row *ent.EventPolicy) *EventPolicy {
	if row == nil {
		return nil
	}

	// Convert conditions from []map[string]any to []PolicyCondition
	var conditions []Condition
	for _, c := range row.Conditions {
		cond := Condition{}
		if f, ok := c["field"].(string); ok {
			cond.Field = f
		}
		if op, ok := c["operator"].(string); ok {
			cond.Operator = op
		}
		cond.Value = c["value"]
		conditions = append(conditions, cond)
	}

	return &EventPolicy{
		ID:                      row.ID,
		Name:                    row.Name,
		Description:             row.Description,
		EventTypes:              row.EventTypes,
		Conditions:              conditions,
		ActionType:              ActionType(row.ActionType.String()),
		MessageTemplate:         row.MessageTemplate,
		RecipientUserIDs:        row.RecipientUserIds,
		RecipientProjectRoleIDs: row.RecipientProjectRoleIds,
		RecipientOrgRoleIDs:     row.RecipientOrgRoleIds,
		RecipientDynamic:        row.RecipientDynamic,
		OrgNodeID:               row.OrgNodeID,
		ProjectID:               row.ProjectID,
		Enabled:                 row.Enabled,
	}
}

// ConditionsToMap converts PolicyConditions to []map[string]any for ent storage
func (e *EventPolicy) ConditionsToMap() []map[string]any {
	result := make([]map[string]any, len(e.Conditions))
	for i, c := range e.Conditions {
		result[i] = map[string]any{
			"field":    c.Field,
			"operator": c.Operator,
			"value":    c.Value,
		}
	}
	return result
}

// MatchesEventType checks if the policy matches a given event type
func (e *EventPolicy) MatchesEventType(eventType string) bool {
	return slices.Contains(e.EventTypes, eventType)
}
