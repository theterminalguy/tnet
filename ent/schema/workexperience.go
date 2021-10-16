package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// WorkExperience holds the schema definition for the WorkExperience entity.
type WorkExperience struct {
	ent.Schema
}

// Mixins for the WorkExperience.
func (WorkExperience) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
		BelongsToMixin{
			ParentName: "applicant",
			ParentType: Applicant.Type,
			Ref:        "work_experiences",
			ForeignKey: "applicant_id",
		},
	}
}

// Fields of the WorkExperience.
func (WorkExperience) Fields() []ent.Field {
	return []ent.Field{
		field.String("company_name"),
		field.String("location"),
		field.String("job_title"),

		field.Text("description"),

		field.Time("start_date"),
		field.Time("end_date").
			Optional(),
	}
}

// Edges of the WorkExperience.
func (WorkExperience) Edges() []ent.Edge {
	return nil
}
