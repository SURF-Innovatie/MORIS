package schema

import (
	"time"

	"entgo.io/contrib/entoas"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type EventPolicy struct {
	ent.Schema
}

func (EventPolicy) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("name"),
		field.String("description").Optional().Nillable(),

		// Trigger: event types + optional conditions
		field.Strings("event_types"),
		// Conditions stored as JSON array of objects: [{field, operator, value}, ...]
		field.JSON("conditions", []map[string]any{}).
			Optional().
			Annotations(entoas.Skip(true)),

		// Action configuration
		field.Enum("action_type").Values("notify", "request_approval"),
		field.String("message_template").Optional().Nillable(),

		// Recipients (at least one should be set)
		field.JSON("recipient_user_ids", []uuid.UUID{}).
			Optional().
			Annotations(entoas.Skip(true)),
		field.JSON("recipient_project_role_ids", []uuid.UUID{}).
			Optional().
			Annotations(entoas.Skip(true)),
		field.JSON("recipient_org_role_ids", []uuid.UUID{}).
			Optional().
			Annotations(entoas.Skip(true)),
		field.Strings("recipient_dynamic").Optional(), // "project_members", "project_owner", "org_admins"

		// Scope: either org_node_id OR project_id is set (not both)
		field.UUID("org_node_id", uuid.UUID{}).Optional().Nillable(),
		field.UUID("project_id", uuid.UUID{}).Optional().Nillable(),

		field.Bool("enabled").Default(true),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (EventPolicy) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("org_node", OrganisationNode.Type).
			Field("org_node_id").
			Unique(),
	}
}

func (EventPolicy) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("org_node_id"),
		index.Fields("project_id"),
	}
}
