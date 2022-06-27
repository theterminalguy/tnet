package schema

import "entgo.io/ent"

// Payment holds the schema definition for the Payment entity.
type Payment struct {
	ent.Schema
}

// Fields of the Payment.
func (Payment) Fields() []ent.Field {
	return nil
}

// Edges of the Payment.
func (Payment) Edges() []ent.Edge {
	return nil
}
