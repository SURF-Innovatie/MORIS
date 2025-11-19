package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type OrganisationChangedEvent struct {
	ent.Schema
}

func (OrganisationChangedEvent) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("organisation_id"),
		field.String("organisation_name"),
	}
}

func (OrganisationChangedEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", Event.Type).
			Ref("organisation_changed").
			Unique().
			Required(),
	}
}
