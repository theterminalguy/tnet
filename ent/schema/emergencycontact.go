package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// EmergencyContact holds the schema definition for the EmergencyContact entity.
type EmergencyContact struct {
	ent.Schema
}

// Mixins for the EmergencyContact.
func (EmergencyContact) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
		BelongsToMixin{
			ParentName: "talent",
			ParentType: Talent.Type,
			Ref:        "emergency_contacts",
			ForeignKey: "talent_id",
		},
	}
}

// Fields of the EmergencyContact.
func (EmergencyContact) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("phone_number"),
		field.String("address"),
		field.String("relationship"),

		field.String("email").
			Unique(),
	}
}

// Edges of the EmergencyContact.
func (EmergencyContact) Edges() []ent.Edge {
	return nil
}
