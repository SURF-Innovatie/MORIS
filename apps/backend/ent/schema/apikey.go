package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// APIKey holds the schema definition for the APIKey entity.
type APIKey struct {
	ent.Schema
}

// Fields of the APIKey.
func (APIKey) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable().
			Unique(),
		field.UUID("user_id", uuid.UUID{}),
		field.String("name").
			NotEmpty(), // User-defined label, e.g., "Power BI"
		field.String("key_hash").
			Sensitive(), // SHA-256 hash of the key
		field.String("key_prefix"), // First 8 chars for identification
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("last_used_at").
			Optional().
			Nillable(),
		field.Time("expires_at").
			Optional().
			Nillable(), // Optional expiration
		field.Bool("is_active").
			Default(true),
	}
}

// Edges of the APIKey.
func (APIKey) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("api_keys").
			Field("user_id").
			Unique().
			Required(),
	}
}

// Indexes of the APIKey.
func (APIKey) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("key_prefix"),
		index.Fields("key_hash").
			Unique(),
	}
}
