package schema

import (
	"time"

	"entgo.io/contrib/entoas"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/SURF-Innovatie/MORIS/internal/types"
	"github.com/google/uuid"
)

// Page holds the schema definition for the Page entity.
type Page struct {
	ent.Schema
}

// Fields of the Page.
func (Page) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("title").NotEmpty(),
		field.String("slug").NotEmpty(),
		field.JSON("content", []types.Section{}).
			Optional().
			Annotations(entoas.Skip(true)),
		field.Enum("type").
			Values("project", "profile").
			Default("project"),
		field.Bool("is_published").Default(false),

		// Polymorphic-style associations
		// Project is event-sourced, so we store ID only
		field.UUID("project_id", uuid.UUID{}).
			Optional().
			Nillable(),

		// User is an Ent entity, so we use an edge (see below)
		// This field is the FK for the edge
		field.UUID("user_id", uuid.UUID{}).
			Optional().
			Nillable(),

		field.Time("created_at").Immutable().Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Page.
func (Page) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner_user", User.Type).
			Ref("pages").
			Field("user_id").
			Unique(),
	}
}

// Indexes of the Page.
func (Page) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("slug").Unique(),
		index.Fields("project_id"), // Index for fast lookups by project
		index.Fields("user_id"),    // Index for fast lookups by user
	}
}
