package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// JobPayment holds the schema definition for the JobPayment entity.
type JobPayment struct {
	ent.Schema
}

func (JobPayment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
		BelongsToMixin{
			ParentName: "job",
			ParentType: Job.Type,
			Ref:        "job_payments",
			ForeignKey: "job_id",
		},
	}
}

// Fields of the JobPayment.
func (JobPayment) Fields() []ent.Field {
	return []ent.Field{
		field.Float("amount").Default(0),
		field.Time("paid_to").Nillable().Optional(),
		field.String("ref_id").Default(""),
		field.String("message"),
		field.String("currency").Default(""),
		field.Text("payment_link"),
		field.JSON("payload", []string{}).Optional(),
	}
}

// Edges of the JobPayment.
func (JobPayment) Edges() []ent.Edge {
	return nil
}
