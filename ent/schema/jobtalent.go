package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/10hourlabs/tentn/oneword"
)

// JobTalent holds the schema definition for the JobTalent entity.
type JobTalent struct {
	ent.Schema
}

func JobTalentStatuses() []string {
	return []string{
		"screening",
		"shortlisted",
		"interviewing",
		"hired",
		"rejected",
	}
}

// Mixins for JobTalent
// func (JobTalent) Mixin() []ent.Mixin {
// 	return []ent.Mixin{
// 		UUIDMixin{},
// 		TimeStampMixin{},
// 		BelongsToMixin{
// 			ParentName: oneword.Applicant,
// 			ParentType: Applicant.Type,
// 			Ref:        oneword.JobTalents,
// 			ForeignKey: oneword.ApplicantID,
// 		},
// 		BelongsToMixin{
// 			ParentName: oneword.Job,
// 			ParentType: Job.Type,
// 			Ref:        oneword.Talents,
// 			ForeignKey: oneword.JobID,
// 		},
// 	}
// }
func (JobTalent) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
		BelongsToMixin{
			ParentName: oneword.Talent,
			ParentType: Talent.Type,
			Ref:        oneword.JobTalents,
			ForeignKey: oneword.TalentID,
		},
		BelongsToMixin{
			ParentName: oneword.Job,
			ParentType: Job.Type,
			Ref:        oneword.JobTalents,
			ForeignKey: oneword.JobID,
		},
	}
}

// Fields of the JobTalent.
func (JobTalent) Fields() []ent.Field {
	return []ent.Field{
		field.String(oneword.ReferralSource),

		field.Enum(oneword.Status).
			Values(JobTalentStatuses()...).
			Default(oneword.Screening),

		field.Text(oneword.Note).
			Optional(),
	}
}

// Edges of the JobTalent.
func (JobTalent) Edges() []ent.Edge {
	return nil
}
