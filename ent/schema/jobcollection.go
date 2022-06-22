package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

func StatusTypes() []string {
	return []string{
		"pending",
		"approved",
	}
}

// JobCollection holds the schema definition for the JobCollection entity.
type JobCollection struct {
	ent.Schema
}

func (JobCollection) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
		BelongsToMixin{
			ParentName: "users",
			ParentType: User.Type,
			Ref:        "job_collections",
			ForeignKey: "recruiter_id",
		},
	}
}

// Fields of the JobCollection.
func (JobCollection) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").
			Values(StatusTypes()...),
		field.String("title"),
	}
}

// Edges of the JobCollection.
func (JobCollection) Edges() []ent.Edge {
	return nil
}
