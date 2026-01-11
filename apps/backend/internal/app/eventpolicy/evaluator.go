package eventpolicy

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"slices"
	"strings"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// Evaluator evaluates policies against events and executes actions
type Evaluator interface {
	// EvaluateAndExecute finds matching policies for an event and executes their actions
	EvaluateAndExecute(ctx context.Context, event events.Event, project *entities.Project) error
	// CheckApprovalRequired checks if any policy requires approval for the event
	CheckApprovalRequired(ctx context.Context, event events.Event, project *entities.Project) (bool, error)
}

type evaluator struct {
	repo               Repository
	closureProvider    OrgClosureProvider
	recipientResolver  RecipientResolver
	notificationSender NotificationSender
}

// NewEvaluator creates a new policy evaluator
func NewEvaluator(
	repo Repository,
	closureProvider OrgClosureProvider,
	recipientResolver RecipientResolver,
	notificationSender NotificationSender,
) Evaluator {
	return &evaluator{
		repo:               repo,
		closureProvider:    closureProvider,
		recipientResolver:  recipientResolver,
		notificationSender: notificationSender,
	}
}

// CheckApprovalRequired checks if any policy requires approval for the event
func (e *evaluator) CheckApprovalRequired(ctx context.Context, event events.Event, project *entities.Project) (bool, error) {
	if project == nil {
		return false, nil
	}

	// 1. Get all applicable policies (project + org hierarchy)
	policies, err := e.getApplicablePolicies(ctx, event.AggregateID(), project.OwningOrgNodeID)
	if err != nil {
		return false, fmt.Errorf("getting applicable policies: %w", err)
	}

	log.Printf("CheckApprovalRequired: Found %d policies for event %s (Project: %s)", len(policies), event.Type(), project.Id)

	// 2. Filter policies that match this event type and pass conditions
	for _, p := range policies {
		if !p.Enabled {
			log.Printf("Policy %s disabled", p.Name)
			continue
		}
		if !p.MatchesEventType(event.Type()) {
			// log.Printf("Policy %s type mismatch (%v vs %s)", p.Name, p.EventTypes, event.Type())
			continue
		}

		matches := e.evaluateConditions(p.Conditions, event, project)
		log.Printf("Policy %s (Action: %s) match result: %v", p.Name, p.ActionType, matches)

		if p.ActionType == entities.ActionTypeRequestApproval && matches {
			log.Printf("Approval required by policy: %s", p.Name)
			return true, nil
		}
	}

	return false, nil
}

// EvaluateAndExecute finds matching policies and executes their actions
func (e *evaluator) EvaluateAndExecute(ctx context.Context, event events.Event, project *entities.Project) error {
	if project == nil {
		return nil
	}

	// 1. Get all applicable policies (project + org hierarchy)
	policies, err := e.getApplicablePolicies(ctx, event.AggregateID(), project.OwningOrgNodeID)
	if err != nil {
		return fmt.Errorf("getting applicable policies: %w", err)
	}

	// 2. Filter policies that match this event type and pass conditions
	matchingPolicies := lo.Filter(policies, func(p entities.EventPolicy, _ int) bool {
		if !p.Enabled {
			return false
		}
		if !p.MatchesEventType(event.Type()) {
			return false
		}
		return e.evaluateConditions(p.Conditions, event, project)
	})

	log.Printf("EvaluateAndExecute: Event %s matches %d policies", event.Type(), len(matchingPolicies))

	// 3. Separate approval and notification policies
	approvalPolicies := lo.Filter(matchingPolicies, func(p entities.EventPolicy, _ int) bool {
		return p.ActionType == entities.ActionTypeRequestApproval
	})
	notificationPolicies := lo.Filter(matchingPolicies, func(p entities.EventPolicy, _ int) bool {
		return p.ActionType == entities.ActionTypeNotify
	})

	log.Printf("EvaluateAndExecute: Found %d approval polices and %d notification policies", len(approvalPolicies), len(notificationPolicies))

	// 4. Execute approval policies first
	approvalSent := false
	for _, policy := range approvalPolicies {
		if err := e.executeAction(ctx, policy, event, project); err != nil {
			log.Printf("policy action error for %s: %v", policy.ID, err)
		} else {
			approvalSent = true
		}
	}

	// 5. Skip notification policies if an approval was already sent
	if approvalSent {
		log.Printf("skipping notification policies - approval already sent for event %s", event.GetID())
		return nil
	}

	// 6. Execute notification policies
	for _, policy := range notificationPolicies {
		if err := e.executeAction(ctx, policy, event, project); err != nil {
			log.Printf("policy action error for %s: %v", policy.ID, err)
		}
	}

	return nil
}

// getApplicablePolicies returns all policies that could apply to a project
func (e *evaluator) getApplicablePolicies(ctx context.Context, projectID uuid.UUID, orgNodeID uuid.UUID) ([]entities.EventPolicy, error) {
	// Get project-level policies
	projectPolicies, err := e.repo.ListForProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Get org hierarchy policies
	ancestorIDs, err := e.closureProvider.GetAncestorIDs(ctx, orgNodeID)
	if err != nil {
		return nil, err
	}

	allOrgIDs := append([]uuid.UUID{orgNodeID}, ancestorIDs...)
	orgPolicies, err := e.repo.ListForOrgNode(ctx, orgNodeID, allOrgIDs)
	if err != nil {
		return nil, err
	}

	return append(projectPolicies, orgPolicies...), nil
}

// evaluateConditions checks if all conditions pass (AND logic)
func (e *evaluator) evaluateConditions(conditions []entities.PolicyCondition, event events.Event, project *entities.Project) bool {
	for _, cond := range conditions {
		if !e.checkCondition(cond, event, project) {
			return false
		}
	}
	return true // All conditions passed (empty conditions = always true)
}

// checkCondition evaluates a single condition
func (e *evaluator) checkCondition(cond entities.PolicyCondition, event events.Event, project *entities.Project) bool {
	value := e.extractValue(cond.Field, event, project)

	switch cond.Operator {
	case entities.OperatorEquals:
		return e.equals(value, cond.Value)
	case entities.OperatorNotEquals:
		return !e.equals(value, cond.Value)
	case entities.OperatorContains:
		return strings.Contains(fmt.Sprint(value), fmt.Sprint(cond.Value))
	case entities.OperatorStartsWith:
		return strings.HasPrefix(fmt.Sprint(value), fmt.Sprint(cond.Value))
	case entities.OperatorGreaterThan:
		return e.toFloat(value) > e.toFloat(cond.Value)
	case entities.OperatorLessThan:
		return e.toFloat(value) < e.toFloat(cond.Value)
	case entities.OperatorIn:
		return e.isIn(value, cond.Value)
	case entities.OperatorNotIn:
		return !e.isIn(value, cond.Value)
	case entities.OperatorExists:
		return value != nil && value != ""
	case entities.OperatorNotExists:
		return value == nil || value == ""
	default:
		log.Printf("unknown operator: %s", cond.Operator)
		return false
	}
}

// extractValue gets a value from event or project by field path
func (e *evaluator) extractValue(fieldPath string, event events.Event, project *entities.Project) any {
	parts := strings.SplitN(fieldPath, ".", 2)
	if len(parts) != 2 {
		return nil
	}

	source, field := parts[0], parts[1]

	switch source {
	case "event":
		return e.getFieldValue(event, field)
	case "project":
		return e.getFieldValue(project, field)
	case "custom_field":
		// Custom field lookup from project's custom fields map
		if project.CustomFields != nil {
			return project.CustomFields[field]
		}
		return nil
	default:
		return nil
	}
}

// getFieldValue uses reflection to get a field value by name
func (e *evaluator) getFieldValue(obj any, fieldName string) any {
	if obj == nil {
		return nil
	}

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}

	// Try exact field name first
	f := v.FieldByName(fieldName)
	if !f.IsValid() {
		// Try case-insensitive match
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			if strings.EqualFold(t.Field(i).Name, fieldName) {
				f = v.Field(i)
				break
			}
		}
	}

	if !f.IsValid() {
		return nil
	}
	return f.Interface()
}

// Helper functions for comparisons
func (e *evaluator) equals(a, b any) bool {
	return fmt.Sprint(a) == fmt.Sprint(b)
}

func (e *evaluator) toFloat(v any) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		var f float64
		fmt.Sscanf(val, "%f", &f)
		return f
	default:
		return 0
	}
}

func (e *evaluator) isIn(value any, collection any) bool {
	switch c := collection.(type) {
	case []any:
		return slices.Contains(c, value)
	case []string:
		return slices.Contains(c, fmt.Sprint(value))
	default:
		return false
	}
}

// executeAction executes the policy's action (notify or request_approval)
func (e *evaluator) executeAction(ctx context.Context, policy entities.EventPolicy, event events.Event, project *entities.Project) error {
	log.Printf("executeAction: Resolving recipients for policy %s", policy.Name)

	// Resolve all recipients
	userIDs, err := e.resolveAllRecipients(ctx, policy, event.AggregateID(), project.OwningOrgNodeID)
	if err != nil {
		return fmt.Errorf("resolving recipients: %w", err)
	}

	if len(userIDs) == 0 {
		log.Printf("executeAction: No recipients found for policy %s", policy.Name)
		return nil // No recipients to notify
	}

	log.Printf("executeAction: Resolved %d recipients for policy %s. Sending %s...", len(userIDs), policy.Name, policy.ActionType)

	// Build message
	message := e.buildMessage(policy, event, project)

	switch policy.ActionType {
	case entities.ActionTypeNotify:
		return e.notificationSender.SendNotification(ctx, userIDs, event.GetID(), message)
	case entities.ActionTypeRequestApproval:
		return e.notificationSender.SendApprovalRequest(ctx, userIDs, event.GetID(), message)
	default:
		return fmt.Errorf("unknown action type: %s", policy.ActionType)
	}
}

// resolveAllRecipients combines all recipient sources into unique user IDs
func (e *evaluator) resolveAllRecipients(ctx context.Context, policy entities.EventPolicy, projectID, orgNodeID uuid.UUID) ([]uuid.UUID, error) {
	userIDSet := make(map[uuid.UUID]bool)

	// Direct "user" IDs (actually person IDs from frontend, need conversion)
	if len(policy.RecipientUserIDs) > 0 {
		userIDs, err := e.recipientResolver.ResolveUsers(ctx, policy.RecipientUserIDs)
		if err != nil {
			log.Printf("error resolving user IDs: %v", err)
		} else {
			for _, uid := range userIDs {
				userIDSet[uid] = true
			}
		}
	}

	// Project role-based recipients
	for _, roleID := range policy.RecipientProjectRoleIDs {
		users, err := e.recipientResolver.ResolveRole(ctx, roleID, projectID)
		if err != nil {
			log.Printf("error resolving project role %s: %v", roleID, err)
			continue
		}
		for _, uid := range users {
			userIDSet[uid] = true
		}
	}

	// Org role-based recipients
	for _, roleID := range policy.RecipientOrgRoleIDs {
		users, err := e.recipientResolver.ResolveOrgRole(ctx, roleID, orgNodeID)
		if err != nil {
			log.Printf("error resolving org role %s: %v", roleID, err)
			continue
		}
		for _, uid := range users {
			userIDSet[uid] = true
		}
	}

	// Dynamic recipients
	for _, dynType := range policy.RecipientDynamic {
		users, err := e.recipientResolver.ResolveDynamic(ctx, dynType, projectID, orgNodeID)
		if err != nil {
			log.Printf("error resolving dynamic %s: %v", dynType, err)
			continue
		}
		for _, uid := range users {
			userIDSet[uid] = true
		}
	}

	return lo.Keys(userIDSet), nil
}

// buildMessage creates the notification message
func (e *evaluator) buildMessage(policy entities.EventPolicy, event events.Event, project *entities.Project) string {
	if policy.MessageTemplate != nil && *policy.MessageTemplate != "" {
		// TODO: implement template substitution (e.g., {{project.title}})
		return *policy.MessageTemplate
	}

	// Default message based on action type
	if policy.ActionType == entities.ActionTypeRequestApproval {
		return fmt.Sprintf("Approval requested for event '%s' on project '%s'", event.Type(), project.Title)
	}
	return fmt.Sprintf("Event '%s' occurred on project '%s'", event.Type(), project.Title)
}
