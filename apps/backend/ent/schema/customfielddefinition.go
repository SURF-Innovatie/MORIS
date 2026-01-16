package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type CustomFieldDefinition struct {
	ent.Schema
}

func (CustomFieldDefinition) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("name").NotEmpty(),
		field.Enum("type").Values("TEXT", "NUMBER", "BOOLEAN", "DATE"),
		field.Enum("category").Values("PROJECT", "PERSON").Default("PROJECT"),
		field.String("description").Optional(),
		field.Bool("required").Default(false),
		field.String("validation_regex").Optional(),
		field.String("example_value").Optional(),
		field.UUID("organisation_node_id", uuid.UUID{}),
	}
}

func (CustomFieldDefinition) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("organisation", OrganisationNode.Type).
			Ref("custom_field_definitions").
			Field("organisation_node_id").
			Unique().
			Required(),
	}
}
