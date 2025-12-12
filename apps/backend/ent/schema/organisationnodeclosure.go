package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type OrganisationNodeClosure struct {
	ent.Schema
}

func (OrganisationNodeClosure) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),

		// ancestor -> descendant
		field.UUID("ancestor_id", uuid.UUID{}),
		field.UUID("descendant_id", uuid.UUID{}),

		// 0 if same node, 1 if parent->child, etc.
		field.Int("depth").NonNegative(),
	}
}

func (OrganisationNodeClosure) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ancestor", OrganisationNode.Type).
			Field("ancestor_id").
			Unique().
			Required(),
		edge.To("descendant", OrganisationNode.Type).
			Field("descendant_id").
			Unique().
			Required(),
	}
}

func (OrganisationNodeClosure) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ancestor_id", "descendant_id").Unique(),
	}
}
