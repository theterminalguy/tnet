package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	"github.com/10hourlabs/tentn/oneword"
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

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String(oneword.FirstName),
		field.String(oneword.LastName),
		field.String("email"),
		field.Enum("role").
			GoType(userrole.Role("")),

		field.Bool("approved").
			Default(false),

		field.String("photo_url"),
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
	return []ent.Edge{
		edge.To("talents", Talent.Type),
		edge.To("slack_app_installs", SlackAppInstall.Type),
		edge.To("jobs", Job.Type),
		edge.To("email_templates", EmailTemplate.Type),
		edge.To("talent_collections", TalentCollection.Type),
	}
}
