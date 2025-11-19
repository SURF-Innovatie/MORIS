package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type PersonAddedEvent struct {
	ent.Schema
}

func (PersonAddedEvent) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("person_id", uuid.UUID{}).
			Default(uuid.New),
		field.String("person_name"),
	}
}

func (PersonAddedEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", Event.Type).
			Ref("person_added").
			Unique().
			Required(),
	}
}
