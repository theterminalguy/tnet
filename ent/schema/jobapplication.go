package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// JobApplication holds the schema definition for the JobApplication entity.
type JobApplication struct {
	ent.Schema
}

func JobApplicationStatuses() []string {
	return []string{
		"screening",
		"shortlisted",
		"interviewing",
		"hired",
		"rejected",
	}
}

// Mixins for JobApplication
func (JobApplication) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
		BelongsToMixin{
			ParentName: "talent",
			ParentType: Talent.Type,
			Ref:        "job_applications",
			ForeignKey: "talent_id",
		},
		BelongsToMixin{
			ParentName: "job",
			ParentType: Job.Type,
			Ref:        "applications",
			ForeignKey: "job_id",
		},
	}
}

// Fields of the JobApplication.
func (JobApplication) Fields() []ent.Field {
	return []ent.Field{
		field.String("referral_source"),

		field.Enum("status").
			Values(JobApplicationStatuses()...).
			Default("screening"),

		field.Text("note").
			Optional(),
	}
}

// Edges of the JobApplication.
func (JobApplication) Edges() []ent.Edge {
	return nil
}
