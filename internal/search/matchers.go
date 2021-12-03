package search

import (
	"fmt"
	"strings"
)

type SearchMatcher string

// List of supported search matchers
var (
	EQ SearchMatcher = "_eq"
	NE SearchMatcher = "_ne"
	GT SearchMatcher = "_gt"
	LT SearchMatcher = "_lt"
	IN SearchMatcher = "_in"

	NEQ   SearchMatcher = "_neq"
	GTE   SearchMatcher = "_gte"
	LTE   SearchMatcher = "_lte"
	NOTIN SearchMatcher = "_notin"
)

// AllowedMatchers is a map of search matchers to ENT's predicates
var AllowedMatchers map[SearchMatcher]string = map[SearchMatcher]string{
	EQ: "EQ",
	NE: "NE",
	GT: "GT",
	LT: "LT",

	NEQ:   "NEQ",
	GTE:   "GTE",
	LTE:   "LTE",
	NOTIN: "NotIn",
}

// ErrInvalidMatcher is returned when a matcher is not allowed
var ErrInvalidMatcher = func(matcher SearchMatcher) error {
	return fmt.Errorf("Invalid matcher: " + string(matcher))
}

type Condition struct {
	Field   string
	Matcher SearchMatcher
	Value   interface{}
}

// Extract takes a map of search conditions and returns a list of conditions
// and a list of fields
func Extract(query map[string]interface{}) (conds []*Condition, errors []error) {
	for k, v := range query {
		lastIndex := strings.LastIndex(k, "_")
		field := NormalizeField(k[:lastIndex])
		matcher := SearchMatcher(k[lastIndex+1:])
		if _, ok := AllowedMatchers[matcher]; !ok {
			errors = append(errors, ErrInvalidMatcher(matcher))
			continue
		}
		c := &Condition{
			Field:   field,
			Matcher: matcher,
			Value:   v,
		}
		conds = append(conds, c)
	}
	return
}

// NormalizeField converts a field to a valid ent field name
func NormalizeField(field string) string {
	return strings.ToLower(strings.Replace(field, ".", "_", -1))
}
