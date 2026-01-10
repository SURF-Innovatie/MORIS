package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type OrganisationRole struct {
	ent.Schema
}

func (OrganisationRole) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("key"), // "admin", "researcher", "custom_role_key"
		field.String("display_name").NotEmpty(),
		field.String("description").Optional(),
		field.UUID("organisation_node_id", uuid.UUID{}),
		field.Strings("permissions").Optional(),
	}
}

func (OrganisationRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("scopes", RoleScope.Type),
		edge.From("organisation", OrganisationNode.Type).
			Ref("organisation_roles").
			Field("organisation_node_id").
			Unique().
			Required(),
	}
}

func (OrganisationRole) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("key", "organisation_node_id").Unique(),
	}
}
