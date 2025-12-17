package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type ProjectRoleAssignedEvent struct {
	ent.Schema
}

func (ProjectRoleAssignedEvent) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("event_id", uuid.UUID{}).Unique(),
		field.UUID("person_id", uuid.UUID{}),
		field.UUID("project_role_id", uuid.UUID{}),
	}
}

func (ProjectRoleAssignedEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", Event.Type).
			Ref("project_role_assigned").
			Field("event_id").
			Unique().
			Required(),

		edge.From("person", Person.Type).
			Ref("project_role_assigned_events").
			Field("person_id").
			Unique().
			Required(),

		edge.From("project_role", ProjectRole.Type).
			Ref("assigned_events").
			Field("project_role_id").
			Unique().
			Required(),
	}
}
func (ProjectRoleAssignedEvent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("event_id").Unique().StorageKey("ux_pra_event"),
		index.Fields("person_id").StorageKey("ix_pra_person"),
		index.Fields("project_role_id").StorageKey("ix_pra_role"),
	}
}
