package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/10hourlabs/tentn/oneword"
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
			ParentName: oneword.Applicant,
			ParentType: Applicant.Type,
			Ref:        oneword.JobApplications,
			ForeignKey: oneword.ApplicantID,
		},
		BelongsToMixin{
			ParentName: oneword.Job,
			ParentType: Job.Type,
			Ref:        oneword.Applications,
			ForeignKey: oneword.JobID,
		},
	}
}

// Fields of the JobApplication.
func (JobApplication) Fields() []ent.Field {
	return []ent.Field{
		field.String(oneword.ReferralSource),

		field.Enum(oneword.Status).
			Values(JobApplicationStatuses()...).
			Default(oneword.Screening),

		field.Text(oneword.Note).
			Optional(),
	}
}

// Edges of the JobApplication.
func (JobApplication) Edges() []ent.Edge {
	return nil
}
