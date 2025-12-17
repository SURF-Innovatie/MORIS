package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Membership struct {
	ent.Schema
}

func (Membership) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),

		field.UUID("person_id", uuid.UUID{}),
		field.UUID("role_scope_id", uuid.UUID{}),
	}
}

func (Membership) Edges() []ent.Edge {
	return []ent.Edge{
		// many memberships -> one person (FK on memberships.person_id)
		edge.To("person", Person.Type).
			Unique().
			Field("person_id").
			Required(),

		// many memberships -> one role_scope (FK on memberships.role_scope_id)
		edge.To("role_scope", RoleScope.Type).
			Unique().
			Field("role_scope_id").
			Required(),
	}
}

func (Membership) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("person_id", "role_scope_id").Unique(),
	}
}
