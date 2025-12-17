package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type OwningOrgNodeChangedEvent struct {
	ent.Schema
}

func (OwningOrgNodeChangedEvent) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("event_id", uuid.UUID{}).
			Unique(),
		field.UUID("owning_org_node_id", uuid.UUID{}),
	}
}

func (OwningOrgNodeChangedEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", Event.Type).
			Ref("owning_org_node_changed").
			Field("event_id").
			Unique().
			Required(),

		edge.To("organisation_node", OrganisationNode.Type).
			Field("owning_org_node_id").
			Unique().
			Required(),
	}
}

func (OwningOrgNodeChangedEvent) Indexes() []ent.Index {
	return []ent.Index{
		// Short names to avoid Postgres truncation collisions
		index.Fields("event_id").Unique().StorageKey("ux_oochg_event"),
		index.Fields("owning_org_node_id").StorageKey("ix_oochg_orgnode"),
	}
}
