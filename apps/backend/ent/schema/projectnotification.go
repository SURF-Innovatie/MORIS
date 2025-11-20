package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type ProjectNotification struct {
	ent.Schema
}

func (ProjectNotification) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.UUID("project_id", uuid.UUID{}),
		field.String("message"),
		field.Time("sent_at").
			Default(time.Now),
	}
}

func (ProjectNotification) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).Unique(),
	}
}
