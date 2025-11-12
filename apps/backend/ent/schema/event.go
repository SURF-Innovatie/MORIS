package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Event struct{ ent.Schema }

func (Event) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("project_id", uuid.UUID{}),
		field.Int("version"),
		field.String("type"),
		field.Bytes("data"),
		field.Bytes("metadata").Optional().Nillable(),
		field.Time("occurred_at").Default(time.Now).Immutable(),
	}
}

func (Event) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("project_id", "version").Unique(),
		index.Fields("project_id", "occurred_at"),
	}
}
