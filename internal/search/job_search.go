package search

import (
	"github.com/10hourlabs/tentn/ent"
	repo "github.com/10hourlabs/tentn/internal/repository"
)

type JobSearch struct {
	JobRepository repo.JobRepository
}

func (*JobSearch) SearchableFields() []string {
	return []string{"uuid", "slug", "title"}
}

func (*JobSearch) AllowedMatchers() []string {
	// TODO: let's define the allowed matchers as a type so we don't have to
	// hardcode them here. It's less error prone and prevent typos.
	return []string{"_eq", "_not_eq"}
}

func (s *JobSearch) Search(query map[string]string) []ent.Job {
	// extract the searchable fields and matchers from the keys of the query
	searchableFields, matchers := s.extractSearchableFieldsAndMatchers(query)
}

func (s *JobSearch) extractSearchableFieldsAndMatchers(query map[string]string) ([]string, []string) {
	var searchableFields []string
	var matchers []string
	for key, value := range query {
		if s.isSearchableField(key) {
			searchableFields = append(searchableFields, key)
			matchers = append(matchers, value)
		}
	}
	return searchableFields, matchers
}

func (s *JobSearch) isSearchableField(key string) bool {
	for _, field := range s.SearchableFields() {
		if field == key {
			return true
		}
	}
	return false
}
