package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
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
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable().Unique(),
		field.UUID("person_id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.String("password").Optional().Sensitive(),
		field.Bool("is_sys_admin").Default(false),
		field.Bool("is_active").Default(true),
		// Zenodo OAuth tokens
		field.String("zenodo_access_token").
			Optional().
			Sensitive(),
		field.String("zenodo_refresh_token").
			Optional().
			Sensitive(),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("notifications", Notification.Type),
	}
}
