package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Partner holds the schema definition for the Partner entity.
type Partner struct {
	ent.Schema
}

// Mixins for Partner
func (Partner) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
	}
}

// Fields of the Partner.
func (Partner) Fields() []ent.Field {
	return []ent.Field{
		field.String("CompanyName"),
		field.String("CompanyLocation"),
		field.String("ContactPersonName"),
		field.String("ContactPersonPhoneNumber"),
		field.String("ContactPersonEmail").
			Unique(),
		field.String("WebsiteUrl").
			Unique(),
	}
}

// Edges of the Partner.
func (Partner) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("missions", Mission.Type),
	}
}
