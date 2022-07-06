package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/10hourlabs/tentn/oneword"
)

// Job holds the schema definition for the Job entity.
type Job struct {
	ent.Schema
}

// Mixins for Job
func (Job) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
		BelongsToMixin{
			ParentName: "user",
			ParentType: User.Type,
			Ref:        "jobs",
			ForeignKey: "user_id",
		},
		BelongsToMixin{
			ParentName: "file_uploads",
			ParentType: FileUpload.Type,
			Ref:        "jobs",
			ForeignKey: "attachment_id",
		},
	}
}

func EmploymentTypes() []string {
	return []string{
		"part_time",
		"full_time",
		"contract",
		"na",
	}
}

func JobCategories() []string {
	return []string{
		"engineering",
		"product_design",
		"sales",
		"marketing",
		"na",
	}
}

// Fields of the Job.
func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.Bool(oneword.Hiring).
			Default(false).
			StructTag(`json:"hiring"`),

		field.String(oneword.Title),

		// Todo, add a default title for this
		field.String(oneword.Slug).
			Unique().
			Immutable(),

		field.String(oneword.Location).
			Default(oneword.RemoteEarth),

		field.String(oneword.Summary),

		// TODO: Add endpoint for returning these values
		field.Enum(oneword.Employment).
			Values(EmploymentTypes()...),

		// TODO: Add endpoint for returning these values
		field.Enum(oneword.Category).
			Values(JobCategories()...),

		field.String(oneword.Thumbnail),

		field.JSON(oneword.WeHave, []string{}),

		field.JSON(oneword.Requirements, []string{}),

		field.JSON(oneword.YouHave, []string{}),

		field.String("timezone"),
	}
}

// Indexes for the Job
func (Job) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(oneword.Title),
		index.Fields(oneword.Category),

		index.Fields(oneword.Slug).
			Unique(),
	}
}

// Edges of the Job.
func (Job) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To(oneword.Applications, JobApplication.Type),
	}
}
