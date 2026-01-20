package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// BudgetLineItem holds the schema definition for the BudgetLineItem entity.
type BudgetLineItem struct {
	ent.Schema
}

// Fields of the BudgetLineItem.
func (BudgetLineItem) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable().
			Unique(),
		field.UUID("budget_id", uuid.UUID{}),
		field.Enum("category").
			Values("personnel", "material", "investment", "travel", "management", "grant", "other"),
		field.String("description").
			NotEmpty(),
		field.Float("budgeted_amount"),
		field.Int("year"),
		field.Enum("funding_source").
			Values("subsidy", "cofinancing_cash", "cofinancing_inkind"),
		field.String("nwo_grant_id").
			Optional().
			Nillable().
			Comment("NWO project/grant ID for linking to NWO Open API"),
	}
}

// Edges of the BudgetLineItem.
func (BudgetLineItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("budget", Budget.Type).
			Ref("line_items").
			Field("budget_id").
			Unique().
			Required(),
		edge.To("actuals", BudgetActual.Type),
	}
}

// Indexes of the BudgetLineItem.
func (BudgetLineItem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("budget_id"),
		index.Fields("category"),
		index.Fields("year"),
		index.Fields("funding_source"),
	}
}
