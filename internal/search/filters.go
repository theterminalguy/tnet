package search

type Filter string

var (
	UUID_EQ Filter = "uuid_eq"
	UUID_NE Filter = "uuid_ne"

	SLUG_EQ Filter = "slug_eq"
	SLUG_NE Filter = "slug_ne"

	TITLE_EQ Filter = "title_eq"
	TITLE_NE Filter = "title_ne"
)
