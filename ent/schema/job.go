package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Job holds the schema definition for the Job entity.
type Job struct {
	ent.Schema
}

// Fields of the Job.
func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("uuid", uuid.UUID{}).
			Default(uuid.New),
	
		field.Time("created_at").
			Default(time.Now),

		field.Time("updated_at").
			Default(time.Now),
		
		field.Time("deleted_at"),

		field.Bool("hiring").
			Default(false),

		field.String("title"),

		field.String("slug"),

		field.String("location").
			Default("Remote, Earth"),

		field.String("summary"),

		field.Enum("employment").
			Values("Part Time", "Full Time", "Contract"),

		field.Enum("category").
			Values("Engineering", "Product & Design"),

		field.String("thumbnail"),

		field.JSON("wehave", []string{}),

		field.JSON("requirements", []string{}),

		field.JSON("youhave", []string{}),
	}
}

// Edges of the Job.
func (Job) Edges() []ent.Edge {
	return nil
}
