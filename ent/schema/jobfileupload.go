package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// FileUpload holds the schema definition for the FileUpload entity.
type JobFileUpload struct {
	ent.Schema
}

func (JobFileUpload) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
	}
}

// Fields of the FileUpload.
func (JobFileUpload) Fields() []ent.Field {
	return []ent.Field{
		field.String("file_url"),
	}
}

// Edges of the FileUpload.
func (JobFileUpload) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("jobs", Job.Type),
	}
}
