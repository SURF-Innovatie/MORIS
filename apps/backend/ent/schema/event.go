package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Event struct {
	ent.Schema
}

func (Event) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.UUID("project_id", uuid.UUID{}),
		field.Int("version"),
		field.String("type"),
		field.Enum("status").
			Values("pending", "approved", "rejected").
			Default("pending"),
		field.UUID("created_by", uuid.UUID{}).
			Optional(), // Optional for now to avoid breaking existing data, or we can set a default if we have a system user
		field.Time("occurred_at").
			Default(time.Now),
	}
}

func (Event) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("project_started", ProjectStartedEvent.Type).Unique(),
		edge.To("title_changed", TitleChangedEvent.Type).Unique(),
		edge.To("description_changed", DescriptionChangedEvent.Type).Unique(),
		edge.To("start_date_changed", StartDateChangedEvent.Type).Unique(),
		edge.To("end_date_changed", EndDateChangedEvent.Type).Unique(),
		edge.To("organisation_changed", OrganisationChangedEvent.Type).Unique(),
		edge.To("person_added", PersonAddedEvent.Type).Unique(),
		edge.To("person_removed", PersonRemovedEvent.Type).Unique(),
		edge.To("product_added", ProductAddedEvent.Type).Unique(),
		edge.To("product_removed", ProductRemovedEvent.Type).Unique(),
	}
}
