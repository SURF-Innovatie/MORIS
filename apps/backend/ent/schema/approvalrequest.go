package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type ApprovalRequest struct {
	ent.Schema
}

func (ApprovalRequest) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),

		field.UUID("event_id", uuid.UUID{}).Unique(),
		field.UUID("project_id", uuid.UUID{}),

		field.Enum("status").
			Values("open", "approved", "rejected").
			Default("open"),

		field.Enum("resolution").
			Values("any_one", "all", "quorum").
			Default("any_one"),

		field.Int("quorum").Optional(),

		field.Time("created_at").Default(time.Now),
		field.Time("closed_at").Optional().Nillable(),
	}
}

func (ApprovalRequest) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("event", Event.Type).
			Field("event_id").
			Unique().
			Required(),

		edge.To("assignees", ApprovalAssignee.Type),
	}
}

func (ApprovalRequest) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("project_id", "status"),
	}
}
