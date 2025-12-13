package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type ProjectRoleAssignedEvent struct {
	ent.Schema
}

func (ProjectRoleAssignedEvent) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("person_id", uuid.UUID{}),
		field.UUID("project_role_id", uuid.UUID{}),
	}
}

func (ProjectRoleAssignedEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", Event.Type).
			Ref("project_role_assigned").
			Unique().
			Required(),

		edge.To("person", Person.Type).
			Field("person_id").
			Unique().
			Required(),

		edge.To("project_role", ProjectRole.Type).
			Field("project_role_id").
			Unique().
			Required(),
	}
}
