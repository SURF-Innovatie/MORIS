package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type ProjectRoleUnassignedEvent struct {
	ent.Schema
}

func (ProjectRoleUnassignedEvent) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("event_id", uuid.UUID{}).Unique(),

		field.UUID("person_id", uuid.UUID{}),
		field.UUID("project_role_id", uuid.UUID{}),
	}
}

func (ProjectRoleUnassignedEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", Event.Type).
			Ref("project_role_unassigned").
			Field("event_id").
			Unique().
			Required(),

		edge.From("person", Person.Type).
			Ref("project_role_unassigned_events").
			Field("person_id").
			Unique().
			Required(),

		edge.From("project_role", ProjectRole.Type).
			Ref("unassigned_events").
			Field("project_role_id").
			Unique().
			Required(),
	}
}

func (ProjectRoleUnassignedEvent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("event_id").Unique().StorageKey("ux_pru_event"),
		index.Fields("person_id").StorageKey("ix_pru_person"),
		index.Fields("project_role_id").StorageKey("ix_pru_role"),
	}
}
