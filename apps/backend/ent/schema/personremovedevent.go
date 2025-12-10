package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type PersonRemovedEvent struct {
	ent.Schema
}

func (PersonRemovedEvent) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("person_id", uuid.UUID{}).
			Default(uuid.New),
	}
}

func (PersonRemovedEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", Event.Type).
			Ref("person_removed").
			Unique().
			Required(),
		edge.To("person", Person.Type).
			Field("person_id").
			Unique().
			Required(),
	}
}
