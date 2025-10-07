package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
			ParentName: "talent_collection",
			ParentType: TalentCollection.Type,
			Ref:        "jobs",
			ForeignKey: "talent_collection_id",
		},
	}
}

func EmploymentTypes() []string {
	return []string{
		"part_time",
		"full_time",
		"contract",
	}
}

func JobCategories() []string {
	return []string{
		"engineering",
		"product_design",
		"sales",
		"marketing",
	}
}

// Fields of the Job.
func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("hiring").
			Default(false).
			StructTag(`json:"hiring"`),

		field.String("title"),

		field.String("ats_job_id").Optional(),

		// Todo, add a default title for this
		field.String("slug").
			Unique().
			Immutable(),

		field.String("location").
			Default("Remote, Earth"),

		field.String("summary"),

		// TODO: Add endpoint for returning these values
		field.Enum("employment").
			Values(EmploymentTypes()...),

		// TODO: Add endpoint for returning these values
		field.Enum("category").
			Values(JobCategories()...),

		field.String("thumbnail"),

		field.JSON("we_have", []string{}),

		field.JSON("requirements", []string{}),

		field.JSON("you_have", []string{}),

		field.String("timezone"),
	}
}

// Indexes for the Job
func (Job) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("title"),
		index.Fields("category"),

		index.Fields("slug").
			Unique(),
	}
}

// Edges of the Job.
func (Job) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("applications", JobApplication.Type),
		edge.To("job_payments", JobPayment.Type),
	}
}
