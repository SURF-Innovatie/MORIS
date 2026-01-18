package events

import (
	"context"
	"fmt"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

// Core Event interface - all events must implement
type Event interface {
	isEvent()
	AggregateID() uuid.UUID
	OccurredAt() time.Time
	GetID() uuid.UUID
	Type() string
	FriendlyName() string
	String() string
	CreatedByID() uuid.UUID
	GetStatus() Status
	SetBase(Base)
}

// Optional behavior interfaces - events implement only if needed
type Applier interface {
	Apply(*entities.Project)
}

type Notifier interface {
	NotificationMessage() string
}

type ApprovalNotifier interface {
	ApprovalMessage(projectTitle string) string
}

type HasRelatedIDs interface {
	RelatedIDs() RelatedIDs
}

// RelatedIDs contains optional references to related entities
type RelatedIDs struct {
	PersonID      *uuid.UUID
	ProductID     *uuid.UUID
	ProjectRoleID *uuid.UUID
	OrgNodeID     *uuid.UUID
}

// EventMeta contains metadata about an event type.
// The function fields allow for context-aware checks that can access the database.
type EventMeta struct {
	Type         string
	FriendlyName string

	// CheckAllowed determines if the actor is allowed to trigger this event.
	// Receives context, event, and ent client for DB access.
	// If nil, defaults to true (all users allowed).
	CheckAllowed func(ctx context.Context, event Event, client *ent.Client) bool
}

// IsAllowed checks if the actor is allowed to trigger this event.
func (m EventMeta) IsAllowed(ctx context.Context, event Event, client *ent.Client) bool {
	if m.CheckAllowed == nil {
		return true // Default: all users allowed
	}
	return m.CheckAllowed(ctx, event, client)
}

// Base struct provides common event fields
type Base struct {
	ID              uuid.UUID `json:"id"`
	ProjectID       uuid.UUID `json:"projectId"`
	FriendlyNameStr string    `json:"friendly_name"`
	At              time.Time `json:"at"`
	CreatedBy       uuid.UUID `json:"createdBy"`
	Status          Status    `json:"status"`
}

func NewBase(projectID, actor uuid.UUID, status Status) Base {
	return Base{
		ProjectID: projectID,
		At:        time.Now().UTC(),
		CreatedBy: actor,
		Status:    status,
	}
}

func (b *Base) AggregateID() uuid.UUID { return b.ProjectID }
func (b *Base) FriendlyName() string   { return b.FriendlyNameStr }
func (b *Base) OccurredAt() time.Time  { return b.At }
func (b *Base) GetID() uuid.UUID       { return b.ID }
func (b *Base) CreatedByID() uuid.UUID { return b.CreatedBy }
func (b *Base) GetStatus() Status      { return b.Status }
func (b *Base) SetBase(base Base)      { *b = base }

// Registry
var (
	eventRegistry = make(map[string]func() Event)
	eventMetas    = make(map[string]EventMeta)
)

// RegisterMeta registers an event factory and its metadata
func RegisterMeta(meta EventMeta, factory func() Event) {
	eventRegistry[meta.Type] = factory
	eventMetas[meta.Type] = meta
}

// Create instantiates an event by type
func Create(eventType string) (Event, error) {
	factory, ok := eventRegistry[eventType]
	if !ok {
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}
	return factory(), nil
}

// GetMeta returns metadata for an event type
func GetMeta(eventType string) EventMeta {
	return eventMetas[eventType]
}

// GetAllMetas returns all registered event metadata
func GetAllMetas() []EventMeta {
	metas := make([]EventMeta, 0, len(eventMetas))
	for _, m := range eventMetas {
		metas = append(metas, m)
	}
	return metas
}
