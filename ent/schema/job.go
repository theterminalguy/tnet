package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Job holds the schema definition for the Job entity.
type Job struct {
	ent.Schema
}

// Mixins for the Job
func (Job) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
	}
}

// Fields of the Job.
func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("hiring").
			Default(false),

		field.String("title"),

		field.String("slug").
			Unique(),

		field.String("location").
			Default("Remote, Earth"),

		field.String("summary"),

		field.Enum("employment").
			Values("part_time", "full_time", "contract"),

		field.Enum("category").
			Values("engineering", "product_design", "sales", "marketing"),

		field.String("thumbnail"),

		field.JSON("wehave", []string{}),

		field.JSON("requirements", []string{}),

		field.JSON("youhave", []string{}),
	}
}

// Indexes for the Job
func (Job) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("title"),

		index.Fields("slug").
			Unique(),
	}
}

// Edges of the Job.
func (Job) Edges() []ent.Edge {
	return nil
}

type JobCategory string

// List of valid Job Category enum
const (
	Sales JobCategory = "Sales"
	Design JobCategory = "Design"
	Product JobCategory = "Product"
	Marketing JobCategory = "Marketing"
	Engineering JobCategory = "Engineering"
)

func (JobCategory) Values (categories []string) {
	for _, c := range []JobCategory{Sales, Design, Product, Marketing, Engineering} {
		categories = append(categories, string(c))
	}
	return
}

type Employment string

// List of valid Employment enum
const (
	PartTime Employment = "Part-Time"
	FullTime Employment = "Full-Time"
	Contract Employment = "Contract"
)

func (Employment) Values (employments []string) {
	for _, e := range []Employment{PartTime, FullTime, Contract} {
		employments = append(employments, string(e))
	}
	return
}
