package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Notification struct {
	ent.Schema
}

func (Notification) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("message"),
		field.Enum("type").
			Values("info", "approval_request", "status_update").
			Default("info"),
		field.Bool("read").Default(false),
		field.Time("sent_at").Default(time.Now),

		field.UUID("user_id", uuid.UUID{}),
		field.UUID("event_id", uuid.UUID{}).Optional().Nillable(),

		// LDN/AS2 fields
		field.String("activity_id").Optional().Nillable().
			Comment("External LDN activity ID (URI)"),
		field.String("activity_type").Optional().Nillable().
			Comment("AS2 activity type (e.g., Announce, Offer, Update)"),
		field.Text("payload").Optional().Nillable().
			Comment("Full AS2 JSON-LD payload (serialized JSON)"),
		field.String("origin_service").Optional().Nillable().
			Comment("Origin service URL"),
		field.String("target_service").Optional().Nillable().
			Comment("Target service URL"),
		field.Enum("direction").
			Values("inbound", "outbound", "internal").
			Default("internal").
			Comment("Notification flow direction"),
		field.Enum("delivery_status").
			Values("pending", "delivered", "failed").
			Default("delivered").
			Comment("Delivery status for outbound notifications"),
	}
}

func (Notification) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("notifications").
			Field("user_id").
			Unique().
			Required(),

		edge.From("event", Event.Type).
			Ref("notifications").
			Field("event_id").
			Unique(),
	}
}
