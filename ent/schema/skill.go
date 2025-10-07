package schema

import (
	"entgo.io/ent"
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
		BelongsToMixin{
			ParentName: "talent",
			ParentType: Talent.Type,
			Ref:        "skills",
			ForeignKey: "talent_id",
		},
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

		field.String("note").Optional(),
	}
}

// Indexes for SKill
func (Skill) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name"),
	}
}

// Edges of the Skill.
func (Skill) Edges() []ent.Edge {
	return nil
}
