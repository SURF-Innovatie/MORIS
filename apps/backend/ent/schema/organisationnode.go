package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type OrganisationNode struct {
	ent.Schema
}

func (OrganisationNode) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("name"),
		field.String("ror_id").Optional().Nillable(),

		// Explicit FK column (nullable => root nodes)
		field.UUID("parent_id", uuid.UUID{}).
			Optional().
			Nillable(),
	}
}

func (OrganisationNode) Edges() []ent.Edge {
	return []ent.Edge{
		// FK lives on the child row, so bind it to parent_id
		edge.From("parent", OrganisationNode.Type).
			Ref("children").
			Unique().
			Field("parent_id"),

		edge.To("children", OrganisationNode.Type),

		edge.To("project_roles", ProjectRole.Type),
		edge.To("custom_field_definitions", CustomFieldDefinition.Type),
	}
}

func (OrganisationNode) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("parent_id"),
	}
}
