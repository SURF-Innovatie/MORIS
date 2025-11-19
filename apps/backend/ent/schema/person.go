package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Person struct {
	ent.Schema
}

func (Person) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("name"),
		field.String("given_name").
			Optional().
			Nillable(),
		field.String("family_name").
			Optional().
			Nillable(),
		field.String("email").
			Optional().
			Nillable(),
	}
}

func (Person) Edges() []ent.Edge {
	return nil
}
