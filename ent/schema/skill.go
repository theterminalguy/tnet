package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Skill holds the schema definition for the Skill entity.
type Skill struct {
	ent.Schema
}

// Mixins for Skill
func (Skill) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
	}
}

// Fields of the Skill.
func (Skill) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),

		field.Float32("years_of_experience").
			Min(1.0),

		field.Bool("preferred").
			Default(false),
		
		field.Text("note"),
		
		field.Int("applicant_id"),
	}
}

// Indexes for SKill
func (Skill) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "applicant_id"),
	}
}

// Edges of the Skill.
func (Skill) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("applicant", Applicant.Type).
			Ref("skills").
			Unique().
			Field("applicant_id"),
	}
}
