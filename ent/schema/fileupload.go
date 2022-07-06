package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// FileUpload holds the schema definition for the FileUpload entity.
type FileUpload struct {
	ent.Schema
}

func (FileUpload) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
	}
}

// Fields of the FileUpload.
func (FileUpload) Fields() []ent.Field {
	return []ent.Field{
		field.String("file_url"),
	}
}

// Edges of the FileUpload.
func (FileUpload) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("jobs", Job.Type),
	}
}
