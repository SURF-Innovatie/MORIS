package schema

import (
	"entgo.io/contrib/entoas"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Person struct {
	ent.Schema
}

func (Person) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).Unique().Immutable(),
		field.UUID("user_id", uuid.UUID{}).
			Default(uuid.New).Unique().Optional(),
		field.String("orcid_id").Optional().Unique(),
		field.String("name"),
		field.String("given_name").
			Optional().
			Nillable(),
		field.String("family_name").
			Optional().
			Nillable(),
		field.String("email").Unique(),
		field.String("avatar_url").Optional().
			Nillable(),
		field.String("description").
			Optional().
			Nillable(),
		field.JSON("org_custom_fields", map[string]interface{}{}).
			Optional().
			Annotations(entoas.Skip(true)),
	}
}

func (Person) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("memberships", Membership.Type).
			Ref("person"),

		edge.To("products", Product.Type),
		edge.To("portfolio", Portfolio.Type).Unique(),
	}
}
