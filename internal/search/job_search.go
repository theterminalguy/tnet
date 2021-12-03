package search

import (
	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/job"
	"github.com/10hourlabs/tentn/ent/predicate"
	repo "github.com/10hourlabs/tentn/internal/repository"
)

type JobSearch struct {
	JobRepository repo.JobRepository
}

func (*JobSearch) SearchableFields() []string {
	// TODO: We should only allow search for indexed fields
	// Some other fields make sense to search, let's add them later
	return []string{"uuid", "slug", "title"}
}

func (*JobSearch) AllowedMatchers() []SearchMatcher {
	// TODO: let's define the allowed matchers as a type so we don't have to
	// hardcode them here. It's less error prone and prevent typos.
	return []SearchMatcher{EQ, NEQ}

}

func (*JobSearch) PossibleFilters() map[string]interface{} {
	m := make(map[string]interface{})
	m["uuid_eq"] = job.UUIDEQ
	m["uuid_neq"] = job.UUIDNEQ

	m["slug_eq"] = job.SlugEQ
	m["slug_neq"] = job.SlugNEQ

	m["title_eq"] = job.TitleEQ
	m["title_neq"] = job.TitleNEQ
	return m
}

func (s *JobSearch) Search(query map[string]string) []ent.Job {
	var ps []predicate.Job

	pf := s.PossibleFilters()
	for key, value := range query {
		if _, ok := pf[key]; !ok {
			continue
		}
		t := pf[key].(func(string) predicate.Job)
		ps = append(ps, t(value))
	}
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
