package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Education holds the schema definition for the Education entity.
type Education struct {
	ent.Schema
}

func (Education) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
		BelongsToMixin{
			ParentName: "applicant",
			ParentType: Applicant.Type,
			Ref:        "educations",
			ForeignKey: "applicant_id",
		},
	}
}

// Fields of the Education.
func (Education) Fields() []ent.Field {
	return []ent.Field{
		field.String("institution_name"),
		field.String("location"),
		field.String("degree"),
		field.String("program"),

		field.Text("overview"),

		field.Time("start_date"),
		field.Time("end_date").
			Optional(),
	}
}

// Edges of the Education.
func (Education) Edges() []ent.Edge {
	return nil
}
