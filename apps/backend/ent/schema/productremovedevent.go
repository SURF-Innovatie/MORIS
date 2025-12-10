package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type ProductRemovedEvent struct {
	ent.Schema
}

func (ProductRemovedEvent) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("product_id", uuid.UUID{}).
			Default(uuid.New),
	}
}

func (ProductRemovedEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", Event.Type).
			Ref("product_removed").
			Unique().
			Required(),
		edge.To("product", Product.Type).
			Field("product_id").
			Unique().
			Required(),
	}
}
