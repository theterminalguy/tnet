package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

// Mixins for Skill
func (Session) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
		BelongsToMixin{
			ParentName: "user",
			ParentType: User.Type,
			Ref:        "sessions",
			ForeignKey: "user_id",
		},
	}
}

// Fields of the Session.
func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.String("session_id").
			Unique(),
		field.String("encoded"),
		field.String("user_agent").Optional(),
		field.String("ip_agent").Optional(),
		field.String("team_id").Optional(),
	}
}

// Edges of the Session.
func (Session) Edges() []ent.Edge {
	return nil
}
