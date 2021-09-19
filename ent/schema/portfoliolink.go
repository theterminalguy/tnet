package schema

import (
	"net/url"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
		field.String("url").
			// TODO extract to an external validation package
			// or use a popular validation library
			Validate(func(s string) error {
				if _, err := url.ParseRequestURI(s); err != nil {
					return err
				}
				return nil
			}),

		field.String("name"),

		field.Int("applicant_id").
			Optional(),
	}
}

// Indexes for the PortfolioLink
func (PortfolioLink) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("applicant_id"),
	}
}

// Edges of the PortfolioLink.
func (PortfolioLink) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("applicant", Applicant.Type).
			Ref("portfoliolinks").
			Unique().
			Field("applicant_id"),
	}
}
