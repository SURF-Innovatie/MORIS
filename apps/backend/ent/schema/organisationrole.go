package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type OrganisationRole struct {
	ent.Schema
}

func (OrganisationRole) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("key").Unique(), // "admin", "researcher", "students"
		field.Bool("has_admin_rights").Default(false),
	}
}

func (OrganisationRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("scopes", RoleScope.Type),
	}
}
