package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type RoleScope struct {
	ent.Schema
}

func (RoleScope) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),

		field.UUID("role_id", uuid.UUID{}),
		field.UUID("root_node_id", uuid.UUID{}),
	}
}

func (RoleScope) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role", OrganisationRole.Type).
			Ref("scopes").
			Field("role_id").
			Unique().
			Required(),

		edge.To("root_node", OrganisationNode.Type).
			Field("root_node_id").
			Unique().
			Required(),

		edge.From("memberships", Membership.Type).
			Ref("role_scope"),
	}
}

func (RoleScope) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("role_id", "root_node_id").Unique(),
	}
}
