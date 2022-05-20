package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// SearchLog holds the schema definition for the SearchLog entity.
type SearchLog struct {
	ent.Schema
}

// Mixins of the SearchLog.
func (SearchLog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeStampMixin{},
	}
}

// Fields of the SearchLog.
func (SearchLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("query"),
		field.Int("result_count"),
		field.String("platform"),
		field.String("platform_user_id"),
		field.String("platform_team_id"),
	}
}

// Edges of the SearchLog.
func (SearchLog) Edges() []ent.Edge {
	return nil
}

// Indexes of the SearchLog.
func (SearchLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("platform", "platform_team_id", "platform_user_id"),
	}
}
