// Package ldn provides types for Linked Data Notifications (LDN) and Activity Streams 2.0 (AS2).
// Implements W3C LDN (https://www.w3.org/TR/ldn/) and COAR Notify Protocol v1.0.1.
package ldn

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DefaultContext is the required @context for COAR Notify notifications.
var DefaultContext = []string{
	"https://www.w3.org/ns/activitystreams",
	"https://coar-notify.net",
}

// Activity represents an AS2 Activity per Activity Streams 2.0 and COAR Notify specs.
// See: https://www.w3.org/TR/activitystreams-core/#activities
type Activity struct {
	Context   []string   `json:"@context"`
	ID        string     `json:"id"`
	Type      []string   `json:"type"`
	Actor     *Actor     `json:"actor,omitempty"`
	Object    *Object    `json:"object"`
	Origin    *Service   `json:"origin"`
	Target    *Service   `json:"target"`
	InReplyTo string     `json:"inReplyTo,omitempty"`
	Summary   string     `json:"summary,omitempty"`
	Published *time.Time `json:"published,omitempty"`
}

// NewActivity creates a new Activity with default context and generated ID.
func NewActivity(activityType string, originURL string) *Activity {
	now := time.Now().UTC()
	return &Activity{
		Context:   DefaultContext,
		ID:        "urn:uuid:" + uuid.New().String(),
		Type:      []string{activityType},
		Origin:    NewService(originURL),
		Published: &now,
	}
}

// Service represents an LDN Service (origin or target) per COAR Notify.
// See: https://coar-notify.net/specification/1.0.1/
type Service struct {
	ID    string `json:"id"`
	Type  string `json:"type"`            // "Service"
	Inbox string `json:"inbox,omitempty"` // LDN inbox URL
	Name  string `json:"name,omitempty"`
}

// NewService creates a Service with the given ID.
func NewService(id string) *Service {
	return &Service{
		ID:   id,
		Type: "Service",
	}
}

// WithInbox sets the inbox URL on a Service.
func (s *Service) WithInbox(inbox string) *Service {
	s.Inbox = inbox
	return s
}

// Actor represents an AS2 Actor (person, organization, application).
// See: https://www.w3.org/TR/activitystreams-vocabulary/#actor-types
type Actor struct {
	ID   string `json:"id"`
	Type string `json:"type"` // Person, Organization, Application, Service
	Name string `json:"name,omitempty"`
}

// NewPerson creates a Person actor.
func NewPerson(id, name string) *Actor {
	return &Actor{ID: id, Type: "Person", Name: name}
}

// NewOrganization creates an Organization actor.
func NewOrganization(id, name string) *Actor {
	return &Actor{ID: id, Type: "Organization", Name: name}
}

// Object represents the focus of an activity.
// See: https://www.w3.org/TR/activitystreams-vocabulary/#object-types
type Object struct {
	ID         string   `json:"id"`
	Type       []string `json:"type,omitempty"`
	IETFCiteAs string   `json:"ietf:cite-as,omitempty"` // Preferred citation URL
	MediaType  string   `json:"mediaType,omitempty"`
	Name       string   `json:"name,omitempty"`
	Content    string   `json:"content,omitempty"`
}

// NewObject creates an Object with an ID.
func NewObject(id string, objectTypes ...string) *Object {
	return &Object{
		ID:   id,
		Type: objectTypes,
	}
}

// Direction indicates the flow direction of a notification.
type Direction string

const (
	DirectionInbound  Direction = "inbound"  // Received from external service
	DirectionOutbound Direction = "outbound" // Sent to external service
	DirectionInternal Direction = "internal" // Internal notification
)

// DeliveryStatus indicates the delivery status of an outbound notification.
type DeliveryStatus string

const (
	DeliveryPending   DeliveryStatus = "pending"
	DeliveryDelivered DeliveryStatus = "delivered"
	DeliveryFailed    DeliveryStatus = "failed"
)

// MarshalJSON implements custom JSON marshaling for Activity.
func (a *Activity) MarshalJSON() ([]byte, error) {
	type Alias Activity
	return json.Marshal((*Alias)(a))
}

// UnmarshalJSON implements custom JSON unmarshaling for Activity.
func (a *Activity) UnmarshalJSON(data []byte) error {
	type Alias Activity
	return json.Unmarshal(data, (*Alias)(a))
}

// Validate checks if the Activity has required fields per COAR Notify spec.
func (a *Activity) Validate() error {
	if len(a.Context) == 0 {
		return ErrMissingContext
	}
	if a.ID == "" {
		return ErrMissingID
	}
	if len(a.Type) == 0 {
		return ErrMissingType
	}
	if a.Origin == nil || a.Origin.ID == "" {
		return ErrMissingOrigin
	}
	if a.Target == nil || a.Target.ID == "" || a.Target.Inbox == "" {
		return ErrMissingTarget
	}
	if a.Object == nil || a.Object.ID == "" {
		return ErrMissingObject
	}
	return nil
}
