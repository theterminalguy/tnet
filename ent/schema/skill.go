package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/10hourlabs/tentn/oneword"
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
			ParentName: oneword.Talent,
			ParentType: Talent.Type,
			Ref:        oneword.Skills,
			ForeignKey: "talent_id",
		},
	}
}

// Fields of the Skill.
func (Skill) Fields() []ent.Field {
	return []ent.Field{
		field.String(oneword.Name),

		field.Float32(oneword.YearsOfExperience).
			Min(1.0),

		field.Bool(oneword.Preferred).
			Default(false),

		field.String(oneword.Note).Optional(),
	}
}

// Indexes for SKill
func (Skill) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(oneword.Name),
	}
}

// Edges of the Skill.
func (Skill) Edges() []ent.Edge {
	return nil
}
