package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Product struct {
	ent.Schema
}

func (Product) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("name"),
		field.String("language").
			Optional().
			Nillable(),
		field.Int("Type").
			Optional(),
		field.String("doi").
			Optional().
			Nillable(),
	}
}

func (Product) Edges() []ent.Edge {
	return nil
}
