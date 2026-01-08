package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type ProjectRole struct {
	ent.Schema
}

func (ProjectRole) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),

		// stable identifier like "contributor", "lead"
		field.String("key").NotEmpty(),

		// optional human label; can equal key
		field.String("name").NotEmpty(),

		// strictly linked to an organisation node
		field.UUID("organisation_node_id", uuid.UUID{}),

		// archived_at marks the role as soft-deleted
		field.Time("archived_at").Optional().Nillable(),
	}
}

func (ProjectRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("organisation", OrganisationNode.Type).
			Ref("project_roles").
			Field("organisation_node_id").
			Unique().
			Required(),
	}
}

func (ProjectRole) Indexes() []ent.Index {
	return []ent.Index{
		// keys must be unique within an organisation
		index.Fields("key", "organisation_node_id").Unique(),
	}
}
