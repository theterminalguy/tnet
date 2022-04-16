package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Oauth2Token holds the schema definition for the Oauth2Token entity.
type Oauth2Token struct {
	ent.Schema
}

// Mixins for Oauth2Token
func (Oauth2Token) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
		BelongsToMixin{
			ParentName: "user",
			ParentType: User.Type,
			Ref:        "oauth2_tokens",
			ForeignKey: "user_id",
		},
		BelongsToMixin{
			ParentName: "oauth2client",
			ParentType: Oauth2Client.Type,
			Ref:        "oauth2_tokens",
			ForeignKey: "oauth2_client_id",
		},
	}
}

// Fields of the Oauth2Token.
func (Oauth2Token) Fields() []ent.Field {
	return []ent.Field{
		// redirect_uri is optional based on the grant type.
		field.String("redirect_uri").
			Optional(),

		field.String("scopes"),
		field.String("code"),
		field.String("code_challenge"),
		field.String("code_challenge_method"),
		field.String("access_token"),
		field.String("refresh_token"),

		field.Time("code_created_at").
			Immutable(),
		field.Time("access_token_created_at").
			Immutable(),
		field.Time("refresh_token_created_at").
			Immutable(),

		field.Int64("code_expires_in"),
		field.Int64("access_token_expires_in"),
		field.Int64("refresh_token_expires_in"),
	}
}

// Indexes for Oauth2Token
func (Oauth2Token) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(
			"code",
			"access_token",
			"refresh_token",
		),
	}
}

// Edges of the Oauth2Token.
func (Oauth2Token) Edges() []ent.Edge {
	return nil
}
