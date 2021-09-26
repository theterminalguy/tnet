package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
)

type TimeStampMixin struct {
	mixin.Schema
}

func (TimeStampMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time(oneword.CreatedAt).
			Immutable().
			Default(time.Now),

		field.Time(oneword.UpdatedAt).
			Default(time.Now).
			UpdateDefault(time.Now),

		field.Time(oneword.DeletedAt).
			Nillable().
			Optional(),
	}
}

type UUIDMixin struct {
	mixin.Schema
}

func (UUIDMixin) Fields() []ent.Field {
	return []ent.Field{
		field.UUID(oneword.UUID, uuid.UUID{}).
			Unique().
			Default(uuid.New),
	}
}

func (UUIDMixin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(oneword.UUID).
			Unique(),
	}
}

// TODO: Blog post on creating resuable mixins
type BelongsToMixin struct {
	ParentName string
	ParentType interface{}
	Ref        string
	ForeignKey string
	mixin.Schema
}

func (b BelongsToMixin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From(b.ParentName, b.ParentType).
			Ref(b.Ref).
			Unique().
			Field(b.ForeignKey),
	}
}

func (b BelongsToMixin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(b.ForeignKey),
	}
}

func (b BelongsToMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Int(b.ForeignKey).
			Optional(),
	}
}
