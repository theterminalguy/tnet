package schema

import (
	"entgo.io/ent/schema/edge"
	"github.com/10hourlabs/tentn/oneword"

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
		field.String(oneword.FirstName),
		field.String(oneword.LastName),
		field.String(oneword.PreferredName),

		field.String(oneword.Pronoun),

		field.String(oneword.PreferredJobTitle),

		field.Bool("is_available"),

		field.Time(oneword.ProfessionalStartDate).StructTag(`json:"career_start_date"`),

		field.String(oneword.Email).
			Unique(),

		// TODO: remove unique constraint
		// due to migration erro
		field.String(oneword.Phone),

		field.String(oneword.CoutryCode).
			MinLen(2).
			MaxLen(2),

		field.String(oneword.City),

		field.Enum("job_preference").Values(JobPreferences()...),

		field.String("timezone"),

		field.String("state"),

		field.String("professional_summary").Optional(),
	}
}

func (Talent) Indexes() []ent.Index {
	// first_name, last_name, email, phone, country_code, city, preferred name, pronoun
	return []ent.Index{
		index.Fields(oneword.Email, oneword.Phone).
			Unique(),
	}
}

// Edges of the Talent.
func (Talent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To(oneword.PortfolioLinks, PortfolioLink.Type),

		edge.To(oneword.Skills, Skill.Type),

		edge.To(oneword.JobApplications, JobApplication.Type),

		edge.To("work_experiences", WorkExperience.Type),

		edge.To("educations", Education.Type),

		edge.To("emergency_contacts", EmergencyContact.Type),

		edge.To("missions", Mission.Type),
	}
}

func JobPreferences() []string {
	return []string{"remote", "onsite", "flexible"}
}
