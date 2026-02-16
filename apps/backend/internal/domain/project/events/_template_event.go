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
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
	"github.com/google/uuid"
)

const TemplateEventType = "project.template"

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
func (TemplateEvent) Type() string { return TemplateEventType }
func (e TemplateEvent) String() string {
	return fmt.Sprintf("Template event: %s", e.SomeStringField)
}

// -----------------------------------------------------------------------------
// Step 3: Implement OPTIONAL behavior interfaces (only what you need)
// -----------------------------------------------------------------------------

// Applier - implement if event modifies project state
func (e *TemplateEvent) Apply(project *project.Project) {
	// Mutate the project based on this event's data
	// Example: project.Title = e.SomeStringField
}

// Notifier - implement if event should notify project members
// All five methods must be implemented together
func (e *TemplateEvent) NotificationTemplate() string {
	return "Something happened: {{event.SomeStringField}}"
}

func (e *TemplateEvent) ApprovalRequestTemplate() string {
	return "Approval needed for action: {{event.SomeStringField}}"
}

func (e *TemplateEvent) ApprovedTemplate() string {
	return "Action '{{event.SomeStringField}}' has been approved."
}

func (e *TemplateEvent) RejectedTemplate() string {
	return "Action '{{event.SomeStringField}}' has been rejected."
}

func (e *TemplateEvent) NotificationVariables() map[string]string {
	return map[string]string{
		"event.SomeStringField": e.SomeStringField,
	}
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
// Step 4: Decision (REQUIRED)
// -----------------------------------------------------------------------------

type TemplateEventInput struct {
	SomeStringField string
	SomeUUIDField   uuid.UUID
}

func DecideTemplateEvent(
	ctx context.Context,
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *project.Project,
	in TemplateEventInput,
	status Status,
) (Event, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if actor == uuid.Nil {
		return nil, errors.New("actor id is required")
	}
	if cur == nil {
		return nil, errors.New("current project is required")
	}

	// Idempotency / invariants go here:
	// if cur.Title == in.SomeStringField { return nil, nil }

	return &TemplateEvent{
		Base:            NewBase(projectID, actor, status),
		SomeStringField: in.SomeStringField,
		SomeUUIDField:   in.SomeUUIDField,
	}, nil
}

// -----------------------------------------------------------------------------
// Step 5: Register the event (REQUIRED - uncomment and customize)
// -----------------------------------------------------------------------------
func init() {
	RegisterMeta(EventMeta{
		Type:         "project.template", // Must match Type() return value
		FriendlyName: "Template Event",   // Human-readable name for UI

		// CheckApproval and CheckNotification removed - now handled via policy evaluator

		// Optional: Define who can trigger this event (nil = everyone)
		CheckAllowed: func(ctx context.Context, event Event, client *ent.Client) bool {
			// Example: client.Project.Query().Where(...)
			return true
		},
	}, func() Event { return &TemplateEvent{} })

	RegisterDecider[TemplateEventInput](TemplateEventType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur any, in TemplateEventInput, status Status) (Event, error) {
			p, ok := cur.(*project.Project)
			if !ok {
				return nil, fmt.Errorf("expected *entities.Project, got %T", cur)
			}
			return DecideTemplateEvent(ctx, projectID, actor, p, in, status)
		},
	)

	RegisterInputType(TemplateEventType, TemplateEventInput{})
}

// Unused import guard (remove when uncommenting init)
var _ = context.Background
