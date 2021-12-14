package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// SlackAppInstall holds the schema definition for the SlackAppInstall entity.
type SlackAppInstall struct {
	ent.Schema
}

// Mixins for SlackAppInstall
func (SlackAppInstall) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
		BelongsToMixin{
			ParentName: "user",
			ParentType: User.Type,
			Ref:        "slack_app_installs",
			ForeignKey: "user_id",
		},
	}
}

// Fields of the SlackAppInstall.
func (SlackAppInstall) Fields() []ent.Field {
	return []ent.Field{
		field.String("team_id"),
		field.String("team_name"),

		field.String("authed_user_id"),
		field.String("authed_user_email"),

		field.String("app_id"),
		field.String("bot_user_id"),

		field.String("access_token"),
		field.String("token_type"),
		field.String("scope"),

		field.Bool("is_enterprise_install"),
	}
}

// Indexes for SlackAppInstall
func (SlackAppInstall) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("team_id").Unique(),
		index.Fields("team_name", "authed_user_email", "authed_user_id"),
	}
}

// Edges of the SlackAppInstall.
func (SlackAppInstall) Edges() []ent.Edge {
	return nil
}
