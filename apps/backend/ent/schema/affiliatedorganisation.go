package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// AffiliatedOrganisation holds the schema definition for the AffiliatedOrganisation entity.
type AffiliatedOrganisation struct {
	ent.Schema
}

// Fields of the AffiliatedOrganisation.
func (AffiliatedOrganisation) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("name").
			NotEmpty(),
		field.String("kvk_number").
			Optional(),
		field.String("ror_id").
			Optional(),
		field.String("vat_number").
			Optional(),
		field.String("city").
			Optional(),
		field.String("country").
			Optional(),
	}
}

// Edges of the AffiliatedOrganisation.
func (AffiliatedOrganisation) Edges() []ent.Edge {
	return nil
}
