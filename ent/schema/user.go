package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Mixins for User
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
	}
}

// TODO: should we implement a value interface instead?
// See: https://entgo.io/docs/schema-fields#enum-fields
func UserRoles() []string {
	return []string{
		"user",
		"partner",
		"admin",
		"superadmin",      // Boss
		"service-account", // Ex Machina
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("email"),
		field.Enum("role").
			Default("user").
			Values(UserRoles()...),
	}
}

// Indexes for User
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email").
			Unique(),
		index.Fields("role"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
