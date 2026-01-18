package eventpolicy

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events/hydrator"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
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
	hydrator           *hydrator.Hydrator
}

// NewEvaluator creates a new policy evaluator
func NewEvaluator(
	repo Repository,
	closureProvider OrgClosureProvider,
	recipientResolver RecipientResolver,
	notificationSender NotificationSender,
	hydrator *hydrator.Hydrator,
) Evaluator {
	return &evaluator{
		repo:               repo,
		closureProvider:    closureProvider,
		recipientResolver:  recipientResolver,
		notificationSender: notificationSender,
		hydrator:           hydrator,
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

	logrus.Infof("CheckApprovalRequired: Found %d policies for event %s (Project: %s)", len(policies), event.Type(), project.Id)

	// 2. Filter policies that match this event type and pass conditions
	for _, p := range policies {
		if !p.Enabled {
			logrus.Infof("Policy %s disabled", p.Name)
			continue
		}
		if !p.MatchesEventType(event.Type()) {
			// logrus.Infof("Policy %s type mismatch (%v vs %s)", p.Name, p.EventTypes, event.Type())
			continue
		}

		matches := e.evaluateConditions(p.Conditions, event, project)
		logrus.Infof("Policy %s (Action: %s) match result: %v", p.Name, p.ActionType, matches)

		if p.ActionType == entities.ActionTypeRequestApproval && matches {
			logrus.Infof("Approval required by policy: %s", p.Name)
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

	logrus.Infof("EvaluateAndExecute: Event %s matches %d policies", event.Type(), len(matchingPolicies))

	// 3. Separate approval and notification policies
	approvalPolicies := lo.Filter(matchingPolicies, func(p entities.EventPolicy, _ int) bool {
		return p.ActionType == entities.ActionTypeRequestApproval
	})
	notificationPolicies := lo.Filter(matchingPolicies, func(p entities.EventPolicy, _ int) bool {
		return p.ActionType == entities.ActionTypeNotify
	})

	logrus.Infof("EvaluateAndExecute: Found %d approval polices and %d notification policies", len(approvalPolicies), len(notificationPolicies))

	// 4. Determine execution strategy based on event status
	status := event.GetStatus()
	logrus.Infof("EvaluateAndExecute: Processing event %s with status %s", event.GetID(), status)

	if status == events.StatusPending {
		// For pending events:
		// 1. Execute approval policies first
		approvalSent := false
		for _, policy := range approvalPolicies {
			if err := e.executeAction(ctx, policy, event, project); err != nil {
				logrus.Infof("policy action error for %s: %v", policy.ID, err)
			} else {
				approvalSent = true
			}
		}

		// 2. Skip notification policies if an approval was already sent (it will be sent on approval)
		if approvalSent {
			logrus.Infof("skipping notification policies - approval already sent for event %s", event.GetID())
			return nil
		}

		// 3. Execute notification policies (no approval required by any policy)
		for _, policy := range notificationPolicies {
			if err := e.executeAction(ctx, policy, event, project); err != nil {
				logrus.Infof("policy action error for %s: %v", policy.ID, err)
			}
		}
	} else if status == events.StatusApproved {
		// For approved events:
		// Only execute notification policies. Approval policies were already handled.
		// We send notifications now because they were likely skipped during the 'pending' phase.
		for _, policy := range notificationPolicies {
			if err := e.executeAction(ctx, policy, event, project); err != nil {
				logrus.Infof("policy action error for %s: %v", policy.ID, err)
			}
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
		logrus.Infof("unknown operator: %s", cond.Operator)
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
	logrus.Infof("executeAction: Resolving recipients for policy %s", policy.Name)

	// Resolve all recipients
	userIDs, err := e.resolveAllRecipients(ctx, policy, event.AggregateID(), project.OwningOrgNodeID)
	if err != nil {
		return fmt.Errorf("resolving recipients: %w", err)
	}

	if len(userIDs) == 0 {
		logrus.Infof("executeAction: No recipients found for policy %s", policy.Name)
		return nil // No recipients to notify
	}

	logrus.Infof("executeAction: Resolved %d recipients for policy %s. Sending %s...", len(userIDs), policy.Name, policy.ActionType)

	// Build message
	message := e.buildMessage(ctx, policy, event, project)

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
			logrus.Infof("error resolving user IDs: %v", err)
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
			logrus.Infof("error resolving project role %s: %v", roleID, err)
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
			logrus.Infof("error resolving org role %s: %v", roleID, err)
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
			logrus.Infof("error resolving dynamic %s: %v", dynType, err)
			continue
		}
		for _, uid := range users {
			userIDSet[uid] = true
		}
	}

	return lo.Keys(userIDSet), nil
}

// buildMessage creates the notification message
func (e *evaluator) buildMessage(ctx context.Context, policy entities.EventPolicy, event events.Event, project *entities.Project) string {
	template := ""
	vars := make(map[string]string)

	// 1. Try event's rich template first
	if rn, ok := event.(events.RichNotifier); ok {
		if policy.ActionType == entities.ActionTypeRequestApproval {
			template = rn.ApprovalRequestTemplate()
		} else {
			template = rn.NotificationTemplate()
		}

		for k, v := range rn.NotificationVariables() {
			vars[k] = v
		}
	}

	// 2. Try policy's custom template (overrides event template if present)
	if policy.MessageTemplate != nil && *policy.MessageTemplate != "" {
		template = *policy.MessageTemplate
	}

	if template != "" {
		// Hydrate event for entity name resolution
		detailed := e.hydrator.HydrateOne(ctx, event)

		// Add standard vars
		vars["project.Title"] = project.Title
		vars["project.Description"] = project.Description

		if detailed.Creator != nil {
			vars["creator.Name"] = detailed.Creator.Name
		}
		if detailed.Product != nil {
			// Assuming Product has a Name/Title field. Using Name based on common patterns.
			// If compilation fails, will adjust.
			vars["product.Name"] = detailed.Product.Name
		}
		if detailed.Person != nil {
			vars["person.Name"] = detailed.Person.Name
		}
		if detailed.OrgNode != nil {
			vars["org_node.Name"] = detailed.OrgNode.Name
		}
		if detailed.ProjectRole != nil {
			vars["role.Name"] = detailed.ProjectRole.Name
		}

		return e.renderTemplate(template, vars)
	}

	// 3. Fall back to basic Notifier
	if n, ok := event.(events.Notifier); ok {
		return n.NotificationMessage()
	}

	// 4. Default fallback
	if policy.ActionType == entities.ActionTypeRequestApproval {
		return fmt.Sprintf("Approval requested for event '%s' on project '%s'", event.FriendlyName(), project.Title)
	}
	return fmt.Sprintf("Event '%s' occurred on project '%s'", event.FriendlyName(), project.Title)
}

func (e *evaluator) renderTemplate(template string, vars map[string]string) string {
	result := template
	for k, v := range vars {
		result = strings.ReplaceAll(result, fmt.Sprintf("{{%s}}", k), v)
	}
	return result
}
