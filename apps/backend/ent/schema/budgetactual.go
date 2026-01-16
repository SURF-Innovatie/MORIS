package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// BudgetActual holds the schema definition for the BudgetActual entity.
type BudgetActual struct {
	ent.Schema
}

// Fields of the BudgetActual.
func (BudgetActual) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable().
			Unique(),
		field.UUID("line_item_id", uuid.UUID{}),
		field.Float("amount"),
		field.String("description").
			Optional(),
		field.Time("recorded_date").
			Default(time.Now),
		field.String("source").
			Default("manual"), // "manual" | "erp_sync"
		field.String("external_ref").
			Optional(), // For ERP reconciliation
	}
}

// Edges of the BudgetActual.
func (BudgetActual) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("line_item", BudgetLineItem.Type).
			Ref("actuals").
			Field("line_item_id").
			Unique().
			Required(),
	}
}

// Indexes of the BudgetActual.
func (BudgetActual) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("line_item_id"),
		index.Fields("recorded_date"),
		index.Fields("source"),
	}
}
