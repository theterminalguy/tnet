package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Payment holds the schema definition for the Payment entity.
type Payment struct {
	ent.Schema
}

func (Payment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
		BelongsToMixin{
			ParentName: "user",
			ParentType: User.Type,
			Ref:        "payments",
			ForeignKey: "user_id",
		},
	}
}

func StatusTypes() []string {
	return []string{
		"paid",
		"not_paid",
	}
}

// Fields of the Payment.
func (Payment) Fields() []ent.Field {
	return []ent.Field{
		field.Float("amount").Optional().Nillable(),
		field.Enum("status").
			Values(StatusTypes()...),
		field.String("ref_id").Optional().Nillable(),
		field.String("message"),
		field.String("currency").Optional().Nillable(),
		field.Text("payment_link").Optional().Nillable(),
		field.UUID("job_collection_id", uuid.UUID{}).Optional().Nillable(), // this should be foreign key to job_collections table
		field.JSON("payload", []string{}).Optional(),
	}
}

// Edges of the Payment.
func (Payment) Edges() []ent.Edge {
	return nil
}
