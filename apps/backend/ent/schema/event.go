package schema

import (
	"time"

	"entgo.io/contrib/entoas"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Event struct {
	ent.Schema
}

func (Event) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.UUID("project_id", uuid.UUID{}),
		field.Int("version"),
		field.String("type"),
		field.Enum("status").
			Values("pending", "approved", "rejected").
			Default("pending"),
		field.UUID("created_by", uuid.UUID{}).
			Optional(),
		field.Time("occurred_at").
			Default(time.Now),
		field.JSON("data", map[string]interface{}{}).
			Default(func() map[string]interface{} { return map[string]interface{}{} }).
			Annotations(entoas.Skip(true)),
	}
}

func (Event) Edges() []ent.Edge {
	return nil
}
