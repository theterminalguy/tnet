package schema

import (
	"entgo.io/ent/schema/edge"
	"github.com/10hourlabs/tentn/oneword"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Applicant holds the schema definition for the Applicant entity.
type Applicant struct {
	ent.Schema
}

// Mixins for Applicant
func (Applicant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
	}
}

// Fields of the Applicant.
func (Applicant) Fields() []ent.Field {
	return []ent.Field{
		field.String(oneword.FirstName),
		field.String(oneword.LastName),
		field.String(oneword.PreferredName),

		field.String(oneword.Pronoun),

		field.String(oneword.PreferredJobTitle),

		// Set on create.
		// This is the applicant referrer database id
		field.Int(oneword.ReferrerID).
			Optional(),

		// Provided by the user
		// Should be validate against existing
		// Applicants TenTNCode
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

		field.Time(oneword.JoinedTenTNAt),
	}
}

func (Applicant) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(oneword.ReferralCode, oneword.ReferrerID),

		index.Fields(oneword.TenTNCode, oneword.Email, oneword.Phone).
			Unique(),
	}
}

// Edges of the Applicant.
func (Applicant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To(oneword.Referees, Applicant.Type).
			From(oneword.Referrer).
			Unique().
			Field(oneword.ReferrerID),

		edge.To(oneword.PortfolioLinks, PortfolioLink.Type),

		edge.To(oneword.Skills, Skill.Type),

		edge.To(oneword.JobApplications, JobApplication.Type),
	}
}
