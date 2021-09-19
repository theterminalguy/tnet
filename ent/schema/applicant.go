package schema

import (
	"crypto/rand"
	"fmt"

	"entgo.io/ent/schema/edge"

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
		field.String("first_name"),
		field.String("last_name"),
		field.String("preferred_name"),

		field.String("pronoun"),

		field.String("preferred_job_title"),

		field.Int("referrer_id").
			Optional(),

		field.String("referral_code").
			Default("NULL").
			Immutable(),

		field.String("tentn_code").
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

		field.Time("professional_start_date"),

		field.String("email").
			Unique(),

		field.String("phone").
			Unique(),

		field.String("country_code").
			MinLen(2).
			MaxLen(2),

		field.String("city"),

		field.Time("joined_tentn_at"),
	}
}

func (Applicant) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("referral_code", "referrer_id"),

		index.Fields("tentn_code", "email", "phone").
			Unique(),
	}
}

// Edges of the Applicant.
func (Applicant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("referees", Applicant.Type).
			From("referrer").
			Unique().
			Field("referrer_id"),
		
		edge.To("portfoliolinks", PortfolioLink.Type),
	}
}
