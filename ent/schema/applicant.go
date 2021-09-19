package schema

import (
	"crypto/rand"
	"fmt"

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

		field.Int(oneword.ReferrerID).
			Optional(),

		field.String(oneword.ReferralCode).
			Default(oneword.NULL).
			Immutable(),

		field.String(oneword.TenTNCode).
			// TODO extract this default logic for ease of testing
			DefaultFunc(func() (code string) {
				b := make([]byte, 5)
				if _, err := rand.Read(b); err != nil {
					// TODO: log error to external service
					// before going live ensure you aren't panicing
					panic(err)
				}
				code = fmt.Sprintf("%X", b)
				return
			}).
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
