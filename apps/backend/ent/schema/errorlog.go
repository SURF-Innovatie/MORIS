package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ErrorLog holds the schema definition for the ErrorLog entity.
type ErrorLog struct {
	ent.Schema
}

// Fields of the ErrorLog.
func (ErrorLog) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable().Unique(),
		field.String("user_id").Optional().Comment("User ID who triggered the error, if authenticated"),
		field.String("http_method").NotEmpty().Comment("HTTP Method (GET, POST, etc.)"),
		field.String("route").NotEmpty().Comment("The route/path accessed"),
		field.Int("status_code").Comment("HTTP Status Code returned"),
		field.Text("error_message").NotEmpty().Comment("The error message"),
		field.Text("stack_trace").Optional().Comment("Stack trace if available"),
		field.Time("timestamp").Default(time.Now).Immutable().Comment("Time when the error occurred"),
	}
}

// Edges of the ErrorLog.
func (ErrorLog) Edges() []ent.Edge {
	return nil
}
