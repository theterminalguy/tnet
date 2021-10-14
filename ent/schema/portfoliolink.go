package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/10hourlabs/tentn/oneword"
)

// PortfolioLink holds the schema definition for the PortfolioLink entity.
type PortfolioLink struct {
	ent.Schema
}

// Mixins for PortfolioLink
func (PortfolioLink) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
	}
}

// Fields of the PortfolioLink.
func (PortfolioLink) Fields() []ent.Field {
	return []ent.Field{
		field.String(oneword.URL),

		field.String(oneword.Name),

		field.Int(oneword.ApplicantID).
			Optional(),
	}
}

// Indexes for the PortfolioLink
func (PortfolioLink) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(oneword.ApplicantID),
	}
}

// Edges of the PortfolioLink.
func (PortfolioLink) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From(oneword.Applicant, Applicant.Type).
			Ref(oneword.PortfolioLinks).
			Unique().
			Field(oneword.ApplicantID),
	}
}
