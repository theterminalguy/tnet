package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// InternalTask holds the schema definition for the InternalTask entity.
type InternalTask struct {
	ent.Schema
}

// Mixins for InternalTask
func (InternalTask) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
	}
}

// Fields of the InternalTask.
func (InternalTask) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("params"),
		field.String("executed_by"), // this is an email address representing the user who executed the task
		field.Bool("succeeded"),
		field.String("error"),
	}
}

// Edges of the InternalTask.
func (InternalTask) Edges() []ent.Edge {
	return nil
}

// Indexes for InternalTask
func (InternalTask) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "executed_by"),
	}
}
