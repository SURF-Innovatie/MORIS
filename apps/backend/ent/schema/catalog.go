package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Catalog holds the schema definition for the Catalog entity.
type Catalog struct {
	ent.Schema
}

// Fields of the Catalog.
func (Catalog) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Immutable(),
		field.String("name").
			NotEmpty().
			Unique().
			Comment("Internal name for the catalog"),
		field.String("description").
			Optional().
			Comment("Description of the catalog"),
		field.Text("rich_description").
			Optional().
			Comment("Rich HTML description from TipTap editor"),
		field.JSON("project_ids", []uuid.UUID{}).
			Optional().
			Comment("List of project IDs included in this catalog"),
		field.String("title").
			NotEmpty().
			Comment("Display title of the catalog"),
		field.String("logo_url").
			Optional().
			Comment("URL to the logo image"),
		field.String("primary_color").
			Optional().
			Comment("Primary branding color (e.g. hex code)"),
		field.String("secondary_color").
			Optional().
			Comment("Secondary branding color (e.g. hex code)"),
		field.String("accent_color").
			Optional().
			Comment("Accent branding color for badges and highlights (e.g. hex code)"),
		field.String("favicon").
			Optional().
			Comment("URL to the favicon"),
		field.String("font_family").
			Optional().
			Comment("Font family to use for the catalog"),
	}
}

// Edges of the Catalog.
func (Catalog) Edges() []ent.Edge {
	return nil
}
