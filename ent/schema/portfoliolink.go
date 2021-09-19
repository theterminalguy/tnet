package schema

import (
	"net/url"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
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
	}
}

// Edges of the PortfolioLink.
func (PortfolioLink) Edges() []ent.Edge {
	return nil
}
