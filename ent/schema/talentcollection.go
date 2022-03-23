package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// TalentCollection holds the schema definition for the TalentCollection entity.
type TalentCollection struct {
	ent.Schema
}

// Mixins for TalentCollection
func (TalentCollection) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
		BelongsToMixin{
			ParentName: "user",
			ParentType: User.Type,
			Ref:        "talent_collections",
			ForeignKey: "user_id",
		},
	}
}

// Fields of the TalentCollection.
func (TalentCollection) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),

		field.JSON("talent_uuids", []string{}),
	}
}

// Edges of the TalentCollection.
func (TalentCollection) Edges() []ent.Edge {
	return nil
}
