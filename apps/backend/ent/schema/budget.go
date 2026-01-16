package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Budget holds the schema definition for the Budget entity.
type Budget struct {
	ent.Schema
}

// Fields of the Budget.
func (Budget) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable().
			Unique(),
		field.UUID("project_id", uuid.UUID{}),
		field.String("title").
			NotEmpty(),
		field.Text("description").
			Optional(),
		field.Enum("status").
			Values("draft", "submitted", "approved", "locked").
			Default("draft"),
		field.Float("total_amount").
			Default(0),
		field.String("currency").
			Default("EUR"),
		field.Int("version").
			Default(1),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Budget.
func (Budget) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("line_items", BudgetLineItem.Type),
	}
}

// Indexes of the Budget.
func (Budget) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("project_id").
			Unique(),
	}
}
