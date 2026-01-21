package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Portfolio struct {
	ent.Schema
}

func (Portfolio) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("person_id", uuid.UUID{}).Unique(),
		field.String("headline").Optional().Nillable(),
		field.String("summary").Optional().Nillable(),
		field.String("website").Optional().Nillable(),
		field.Bool("show_email").Default(true),
		field.Bool("show_orcid").Default(true),
		field.JSON("pinned_project_ids", []uuid.UUID{}).Optional(),
		field.JSON("pinned_product_ids", []uuid.UUID{}).Optional(),
	}
}

func (Portfolio) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("person", Person.Type).
			Ref("portfolio").
			Field("person_id").
			Unique().
			Required(),
	}
}
