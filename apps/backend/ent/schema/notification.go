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
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("message"),
		field.Enum("type").
			Values("info", "approval_request", "status_update").
			Default("info"),
		field.Bool("read").Default(false),
		field.Time("sent_at").
			Default(time.Now),
	}
}

func (Notification) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).Unique(),
		edge.To("event", Event.Type).Unique(),
	}
}
