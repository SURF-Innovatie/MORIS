package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type ProjectStartedEvent struct {
	ent.Schema
}

func (ProjectStartedEvent) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("title"),
		field.String("description"),
		field.Time("start_date"),
		field.Time("end_date"),
		field.UUID("organisation_id", uuid.UUID{}),
	}
}

func (ProjectStartedEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", Event.Type).
			Ref("project_started").
			Unique().
			Required(),
	}
}
