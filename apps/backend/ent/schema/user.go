package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique(),
		field.String("name").NotEmpty(),
		field.String("email").Unique(),
		field.String("password").NotEmpty().Sensitive(), // Sensitive prevents it from being logged
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
