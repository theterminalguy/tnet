package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// SlackAppUser holds the schema definition for the SlackAppUser entity.
type SlackAppUser struct {
	ent.Schema
}

// Mixins for SlackAppUser
func (SlackAppUser) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
		BelongsToMixin{
			ParentName: "slack_app_install",
			ParentType: SlackAppInstall.Type,
			Ref:        "slack_app_users",
			ForeignKey: "slack_app_install_id",
		},
	}
}

// Fields of the SlackAppUser.
func (SlackAppUser) Fields() []ent.Field {
	return []ent.Field{
		field.String("full_name"),
		field.String("title").Optional(),
		field.String("email"),
		field.String("photo_url"),
		field.String("slack_user_id"),
		field.String("slack_team_id"),
		field.String("timezone"),
		field.String("timezone_label"),
		field.String("locale").Default("en-US"),
		field.Bool("is_bot_user"),
	}
}

// Indexes for SlackAppUser
func (SlackAppUser) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("slack_user_id", "email").Unique(),
		index.Fields("slack_team_id"),
	}
}

// Edges of the SlackAppUser.
func (SlackAppUser) Edges() []ent.Edge {
	return nil
}
