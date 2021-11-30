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

		// Set on create.
		// This is the Talent referrer database id
		field.Int(oneword.ReferrerID).
			Optional().
			StructTag(`json:"-"`),

		// Provided by the user
		// Should be validate against existing
		// Talents TenTNCode
		// This should be optional
		field.String(oneword.ReferralCode).
			Default(oneword.NULL).
			Optional().
			Immutable(),

		// Set on create
		field.String(oneword.TenTNCode).
			Unique().
			Immutable(),

		field.Time(oneword.ProfessionalStartDate),

		field.String(oneword.Email).
			Unique(),

		field.String(oneword.Phone).
			Unique(),

		field.String(oneword.CoutryCode).
			MinLen(2).
			MaxLen(2),

		field.String(oneword.City),

		field.Time(oneword.JoinedTenTNAt).
			Nillable().
			Optional(),
	}
}

func (Talent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(oneword.ReferralCode, oneword.ReferrerID),

		index.Fields(oneword.TenTNCode, oneword.Email, oneword.Phone).
			Unique(),
	}
}

// Edges of the Talent.
func (Talent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To(oneword.Referees, Talent.Type).
			From(oneword.Referrer).
			Unique().
			Field(oneword.ReferrerID),

		edge.To(oneword.PortfolioLinks, PortfolioLink.Type),

		edge.To(oneword.Skills, Skill.Type),

		edge.To(oneword.JobTalents, JobTalent.Type),

		edge.To("work_experiences", WorkExperience.Type),

		edge.To("educations", Education.Type),

		edge.To("emergency_contacts", EmergencyContact.Type),

		edge.To("missions", Mission.Type),
	}
}
