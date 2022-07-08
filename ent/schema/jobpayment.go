package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Payment holds the schema definition for the JobPayment entity.
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

func StatusTypes() []string {
	return []string{
		"paid",
		"not_paid",
	}
}

// Fields of the JobPayment.
func (JobPayment) Fields() []ent.Field {
	return []ent.Field{
		field.Float("amount").Default(0),
		field.Enum("status").
			Values(StatusTypes()...),
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
