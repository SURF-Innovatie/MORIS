package notification

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	entnotif "github.com/SURF-Innovatie/MORIS/ent/notification"
	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationInfo            NotificationType = "info"
	NotificationApprovalRequest NotificationType = "approval_request"
	NotificationStatusUpdate    NotificationType = "status_update"
)

// Direction indicates the flow direction of a notification.
type Direction string

const (
	DirectionInbound  Direction = "inbound"
	DirectionOutbound Direction = "outbound"
	DirectionInternal Direction = "internal"
)

// DeliveryStatus indicates the delivery status of an outbound notification.
type DeliveryStatus string

const (
	DeliveryPending   DeliveryStatus = "pending"
	DeliveryDelivered DeliveryStatus = "delivered"
	DeliveryFailed    DeliveryStatus = "failed"
)

type Notification struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	EventID           *uuid.UUID // optional
	ProjectID         *uuid.UUID // optional, derived from event
	Message           string
	Type              NotificationType
	Read              bool
	SentAt            time.Time
	EventFriendlyName *string

	// LDN/AS2 fields
	ActivityID     *string
	ActivityType   *string
	Payload        *string // Serialized JSON-LD
	OriginService  *string
	TargetService  *string
	Direction      Direction
	DeliveryStatus DeliveryStatus
}

func (n *Notification) FromEnt(row *ent.Notification) *Notification {
	out := &Notification{
		ID:             row.ID,
		Message:        row.Message,
		Type:           NotificationType(row.Type.String()),
		Read:           row.Read,
		SentAt:         row.SentAt,
		UserID:         row.UserID,
		Direction:      Direction(row.Direction.String()),
		DeliveryStatus: DeliveryStatus(row.DeliveryStatus.String()),
	}
	if row.EventID != nil {
		out.EventID = row.EventID
	}
	if row.Edges.Event != nil {
		out.ProjectID = &row.Edges.Event.ProjectID
	}
	out.ActivityID = row.ActivityID
	out.ActivityType = row.ActivityType
	out.Payload = row.Payload
	out.OriginService = row.OriginService
	out.TargetService = row.TargetService
	return out
}

// DirectionFromEnt converts Ent direction enum to domain type.
func DirectionFromEnt(d entnotif.Direction) Direction {
	return Direction(d.String())
}

// DeliveryStatusFromEnt converts Ent delivery status enum to domain type.
func DeliveryStatusFromEnt(s entnotif.DeliveryStatus) DeliveryStatus {
	return DeliveryStatus(s.String())
}
