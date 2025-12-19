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
		field.UUID("project_id", uuid.UUID{}).Default(uuid.New),

		field.String("key").NotEmpty(),
		field.String("name").NotEmpty(),
	}
}

func (ProjectRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("assigned_events", ProjectRoleAssignedEvent.Type),
		edge.To("unassigned_events", ProjectRoleUnassignedEvent.Type),
	}
}

func (ProjectRole) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("project_id", "key").Unique(),
	}
}
