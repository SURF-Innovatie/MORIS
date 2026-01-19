package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// RaidInfo holds the schema definition for the RaidInfo entity.
type RaidInfo struct {
	ent.Schema
}

// Fields of the RaidInfo.
func (RaidInfo) Fields() []ent.Field {
	return []ent.Field{
		field.String("raid_id").Unique().NotEmpty(),
		field.String("schema_uri").Default("https://raid.org/"),
		field.String("registration_agency_id"),
		field.String("registration_agency_schema_uri"),
		field.String("owner_id"),
		field.String("owner_schema_uri"),
		field.Int64("owner_service_point").Optional().Nillable(),

		field.UUID("project_id", uuid.UUID{}),

		field.String("license"),
		field.Int("version").Default(1),
		field.Time("latest_sync").Optional().Nillable(),
		field.Bool("dirty").Default(false),
		field.String("checksum").Optional().Nillable(),
	}
}

// Edges of the RaidInfo.
func (RaidInfo) Edges() []ent.Edge {
	return []ent.Edge{}
}
