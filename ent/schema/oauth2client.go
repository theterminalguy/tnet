package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Oauth2Client holds the schema definition for the Oauth2Client entity.
type Oauth2Client struct {
	ent.Schema
}

// Mixins for Oauth2Client
func (Oauth2Client) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
		BelongsToMixin{
			ParentName: "user",
			ParentType: User.Type,
			Ref:        "oauth2_clients",
			ForeignKey: "user_id",
		},
	}
}

// Fields of the Oauth2Client.
func (Oauth2Client) Fields() []ent.Field {
	return []ent.Field{
		field.String("app_name"),
		field.String("app_description"),
		field.String("app_logo_uri"),
		field.String("app_homepage_uri"),
		field.String("app_privacy_policy_uri"),

		field.JSON("scopes", []string{}),
		field.String("hashed_secret"),
		field.JSON("redirect_uris", []string{}),

		// https://datatracker.ietf.org/doc/html/rfc6749#section-2.1
		field.Enum("client_type").
			Values(ClientTypes()...),

		field.Bool("is_internal").Default(false),
		field.Bool("approved").Default(false),
	}
}

// Edges of the Oauth2Client.
func (Oauth2Client) Edges() []ent.Edge {
	return nil
}

func ClientTypes() []string {
	return []string{
		"confidential",
		"public",
	}
}
