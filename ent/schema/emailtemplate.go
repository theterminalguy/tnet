package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// EmailTemplate holds the schema definition for the EmailTemplate entity.
type EmailTemplate struct {
	ent.Schema
}

func (EmailTemplate) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
		BelongsToMixin{
			ParentName: "user",
			ParentType: User.Type,
			Ref:        "email_templates",
			ForeignKey: "user_id",
		},
	}
}

// Fields of the EmailTemplate.
func (EmailTemplate) Fields() []ent.Field {
	return []ent.Field{
		field.String("from"),
		field.String("subject"),

		field.Text("body"),

		field.Enum("status").
			Values(JobApplicationStatuses()...).
			Default("screening"),

		field.JSON("cc", []string{}),
		field.JSON("bcc", []string{}),
	}
}

// Edges of the EmailTemplate.
func (EmailTemplate) Edges() []ent.Edge {
	return nil
}
