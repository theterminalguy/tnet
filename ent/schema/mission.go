package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Mission holds the schema definition for the Mission entity.
type Mission struct {
	ent.Schema
}

func MissionTypes() []string {
	return []string{
		"internal",
		"external",
	}
}

// Mixins for Mission
func (Mission) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
		BelongsToMixin{
			ParentName: "talent",
			ParentType: Talent.Type,
			Ref:        "missions",
			ForeignKey: "talent_id",
		},
		BelongsToMixin{
			ParentName: "partner",
			ParentType: Partner.Type,
			Ref:        "missions",
			ForeignKey: "partner_id",
		},
	}
}

// Fields of the Mission.
func (Mission) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("mission_type").
			Values(MissionTypes()...),

		field.Time("start_date"),
		field.Time("end_date").
			Optional(),
	}
}

// Edges of the Mission.
func (Mission) Edges() []ent.Edge {
	return nil
}
