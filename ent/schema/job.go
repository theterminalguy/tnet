package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Job holds the schema definition for the Job entity.
type Job struct {
	ent.Schema
}

// Mixins for the Job
func (Job) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
	}
}

// Fields of the Job.
func (Job) Fields() []ent.Field {
	return []ent.Field{
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
