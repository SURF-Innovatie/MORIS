package events

import (
	"regexp"
	"strings"
)

var templateVarRegex = regexp.MustCompile(`\{\{([^}]+)\}\}`)

// ResolveTemplate resolves {{variable}} placeholders in a template string.
// It uses the provided variables map for substitution.
// Variables from event.NotificationVariables() should be passed directly.
// For related entities, use AddDetailedEventVariables to extend the variables map.
func ResolveTemplate(template string, variables map[string]string) string {
	if template == "" {
		return ""
	}

	return templateVarRegex.ReplaceAllStringFunc(template, func(match string) string {
		// Extract variable name from {{variable}}
		varName := strings.TrimPrefix(strings.TrimSuffix(match, "}}"), "{{")
		varName = strings.TrimSpace(varName)

		if val, ok := variables[varName]; ok {
			return val
		}
		// Keep the placeholder if no value found
		return match
	})
}

// AddDetailedEventVariables adds variables from DetailedEvent related entities.
// This extends the base variables from NotificationVariables() with hydrated entity data.
func AddDetailedEventVariables(vars map[string]string, de DetailedEvent) map[string]string {
	if vars == nil {
		vars = make(map[string]string)
	}

	// Add person variables
	if de.Person != nil {
		vars["person.Name"] = de.Person.Name
		vars["person.Email"] = de.Person.Email
	}

	// Add product variables
	if de.Product != nil {
		vars["product.Name"] = de.Product.Name
	}

	// Add project role variables
	if de.ProjectRole != nil {
		vars["project_role.Name"] = de.ProjectRole.Name
	}

	// Add org node variables
	if de.OrgNode != nil {
		vars["org_node.Name"] = de.OrgNode.Name
	}

	// Add creator variables
	if de.Creator != nil {
		vars["creator.Name"] = de.Creator.Name
		vars["creator.Email"] = de.Creator.Email
	}

	return vars
}

// BuildNotificationMessage builds a notification message from a Notifier event.
// It combines the event's template with variables from both NotificationVariables()
// and hydrated DetailedEvent data.
func BuildNotificationMessage(n Notifier, de DetailedEvent) string {
	template := n.NotificationTemplate()
	if template == "" {
		return ""
	}

	vars := n.NotificationVariables()
	if vars == nil {
		vars = make(map[string]string)
	}

	vars = AddDetailedEventVariables(vars, de)
	return ResolveTemplate(template, vars)
}

// BuildApprovalRequestMessage builds an approval request message from a Notifier event.
func BuildApprovalRequestMessage(n Notifier, de DetailedEvent) string {
	template := n.ApprovalRequestTemplate()
	if template == "" {
		return ""
	}

	vars := n.NotificationVariables()
	if vars == nil {
		vars = make(map[string]string)
	}

	vars = AddDetailedEventVariables(vars, de)
	return ResolveTemplate(template, vars)
}

// BuildApprovedMessage builds an approved status message from a Notifier event.
func BuildApprovedMessage(n Notifier, de DetailedEvent) string {
	template := n.ApprovedTemplate()
	if template == "" {
		return ""
	}

	vars := n.NotificationVariables()
	if vars == nil {
		vars = make(map[string]string)
	}

	vars = AddDetailedEventVariables(vars, de)
	return ResolveTemplate(template, vars)
}

// BuildRejectedMessage builds a rejected status message from a Notifier event.
func BuildRejectedMessage(n Notifier, de DetailedEvent) string {
	template := n.RejectedTemplate()
	if template == "" {
		return ""
	}

	vars := n.NotificationVariables()
	if vars == nil {
		vars = make(map[string]string)
	}

	vars = AddDetailedEventVariables(vars, de)
	return ResolveTemplate(template, vars)
}
