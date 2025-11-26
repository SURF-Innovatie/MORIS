package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type ProductAddedEvent struct {
	ent.Schema
}

func (ProductAddedEvent) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("product_id", uuid.UUID{}).
			Default(uuid.New),
	}
}

func (ProductAddedEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", Event.Type).
			Ref("product_added").
			Unique().
			Required(),
	}
}
