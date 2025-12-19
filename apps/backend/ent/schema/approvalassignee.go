package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type ApprovalAssignee struct{ ent.Schema }

func (ApprovalAssignee) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),

		field.UUID("approval_request_id", uuid.UUID{}),
		field.UUID("person_id", uuid.UUID{}),

		field.Enum("source_kind").
			Values("project_role", "org_role"),

		field.String("source_role_key").NotEmpty(),

		field.UUID("source_scope_root_node_id", uuid.UUID{}).
			Optional().
			Nillable(),

		field.Enum("state").
			Values("pending", "approved", "rejected").
			Default("pending"),

		field.Time("decided_at").Optional().Nillable(),
	}
}

func (ApprovalAssignee) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("approval_request", ApprovalRequest.Type).
			Field("approval_request_id").
			Unique().
			Required(),

		edge.To("person", Person.Type).
			Field("person_id").
			Unique().
			Required(),
	}
}

func (ApprovalAssignee) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("approval_request_id", "person_id").Unique(),

		index.Fields("person_id", "state"),
	}
}
