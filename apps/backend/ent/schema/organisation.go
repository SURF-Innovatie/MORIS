package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Organisation struct {
	ent.Schema
}

func (Organisation) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("name"),
	}
}

func (Organisation) Edges() []ent.Edge {
	return nil
}
