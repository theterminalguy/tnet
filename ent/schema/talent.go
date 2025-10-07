package schema

import (
	"entgo.io/ent/schema/edge"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Talent holds the schema definition for the Talent entity.
type Talent struct {
	ent.Schema
}

// Mixins for Talent
func (Talent) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
		BelongsToMixin{
			ParentName: "user",
			ParentType: User.Type,
			Ref:        "talents",
			ForeignKey: "user_id",
		},
	}
}

// Fields of the Talent.
func (Talent) Fields() []ent.Field {
	return []ent.Field{
		field.String("first_name"),
		field.String("last_name"),
		field.String("preferred_name"),

		field.String("pronoun"),

		field.String("preferred_job_title"),

		field.Bool("is_available"),

		field.String("slug").
			Optional(),

		field.Time("professional_start_date").StructTag(`json:"career_start_date"`),

		field.String("email").
			Unique(),

		// TODO: remove unique constraint
		// due to migration erro
		field.String("phone"),

		field.String("country_code").
			MinLen(2).
			MaxLen(2),

		field.String("city"),

		field.Enum("job_preference").Values(JobPreferences()...),

		field.String("timezone"),

		field.String("locale").Default("en-US"),

		field.String("state"),

		field.String("professional_summary").Optional(),
	}
}

func (Talent) Indexes() []ent.Index {
	// first_name, last_name, email, phone, country_code, city, preferred name, pronoun
	return []ent.Index{
		index.Fields("email", "phone").
			Unique(),
	}
}

// Edges of the Talent.
func (Talent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("portfoliolinks", PortfolioLink.Type),

		edge.To("skills", Skill.Type),

		edge.To("job_applications", JobApplication.Type),

		edge.To("work_experiences", WorkExperience.Type),

		edge.To("educations", Education.Type),

		edge.To("emergency_contacts", EmergencyContact.Type),

		edge.To("missions", Mission.Type),
	}
}

func JobPreferences() []string {
	return []string{"remote", "onsite", "flexible"}
}
