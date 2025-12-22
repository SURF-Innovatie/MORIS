package events

// =============================================================================
// TEMPLATE EVENT - Copy this file when creating a new event
// =============================================================================
//
// To create a new event:
// 1. Copy this file and rename it (e.g., my_new_event.go)
// 2. Rename the struct and update the type constant
// 3. Implement only the interfaces your event needs
// 4. Uncomment the init() function and update the RegisterMeta call
// =============================================================================

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

// -----------------------------------------------------------------------------
// Step 1: Define your event struct
// -----------------------------------------------------------------------------
// Embed Base to get common fields (ID, ProjectID, At, CreatedBy, Status).
type TemplateEvent struct {
	Base // Always embed Base

	// Add your event-specific fields here with JSON tags
	SomeStringField string    `json:"some_string_field"`
	SomeUUIDField   uuid.UUID `json:"some_uuid_field"`
}

// -----------------------------------------------------------------------------
// Step 2: Implement the core Event interface (REQUIRED)
// -----------------------------------------------------------------------------

func (TemplateEvent) isEvent()     {}
func (TemplateEvent) Type() string { return "project.template" }
func (e TemplateEvent) String() string {
	return fmt.Sprintf("Template event: %s", e.SomeStringField)
}

// -----------------------------------------------------------------------------
// Step 3: Implement OPTIONAL behavior interfaces (only what you need)
// -----------------------------------------------------------------------------

// Applier - implement if event modifies project state
func (e *TemplateEvent) Apply(project *entities.Project) {
	// Mutate the project based on this event's data
	// Example: project.Title = e.SomeStringField
}

// Notifier - implement if event should notify project members
func (e *TemplateEvent) NotificationMessage() string {
	return fmt.Sprintf("Something happened: %s", e.SomeStringField)
}

// ApprovalNotifier - implement if event requires approval workflow
func (e *TemplateEvent) ApprovalMessage(projectTitle string) string {
	return fmt.Sprintf("Approval needed for action in project '%s'.", projectTitle)
}

// HasRelatedIDs - implement if event references related entities
func (e *TemplateEvent) RelatedIDs() RelatedIDs {
	return RelatedIDs{
		// Set only the fields relevant to your event
		// PersonID:      &e.SomePersonID,
		// ProductID:     &e.SomeProductID,
	}
}

// -----------------------------------------------------------------------------
// Step 4: Register the event (REQUIRED - uncomment and customize)
// -----------------------------------------------------------------------------
//
// func init() {
// 	RegisterMeta(EventMeta{
// 		Type:         "project.template",  // Must match Type() return value
// 		FriendlyName: "Template Event",    // Human-readable name for UI
//
// 		// Optional: Define when approval is required (nil = never)
// 		CheckApproval: func(ctx context.Context, event Event, client *ent.Client) bool {
// 			return false // or add custom logic
// 		},
//
// 		// Optional: Define when notifications should be sent (nil = never)
// 		CheckNotification: func(ctx context.Context, event Event, client *ent.Client) bool {
// 			return true // or add custom logic
// 		},
//
// 		// Optional: Define who can trigger this event (nil = everyone)
// 		CheckAllowed: func(ctx context.Context, event Event, client *ent.Client) bool {
// 			// Example: client.Project.Query().Where(...)
// 			return true
// 		},
// 	}, func() Event { return &TemplateEvent{} })
// }

// Unused import guard (remove when uncommenting init)
var _ = context.Background
