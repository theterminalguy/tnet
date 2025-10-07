package search

import (
	"fmt"
	"net/url"

	"github.com/theterminalguy/tentn/ent"
	"github.com/theterminalguy/tentn/ent/jobapplication"
	"github.com/theterminalguy/tentn/ent/predicate"
	repo "github.com/theterminalguy/tentn/internal/repository"
)

type JobApplicationSearch struct {
	JobApplicationRepository repo.JobApplicationRepository
}

func (*JobApplicationSearch) PossibleFilters() []Filter {
	return []Filter{
		STATUS_EQ,
		STATUS_NEQ,
	}
}

func (s *JobApplicationSearch) Search(qs string) ([]*ent.JobApplication, []error) {
	// TODO: this implementation of search and filters is not reusable but works really well for now
	// and makes it easy for us to decide what can be searchable and what can't.
	// We should probably define a Searchable interface and implement it for each searchable entity
	var ps []predicate.JobApplication
	var errors []error
	query, err := url.ParseQuery(qs)
	if err != nil {
		errors = append(errors, err)
	}

	pf := s.PossibleFilters()
	for _, filter := range pf {
		f := string(filter)
		if vv, ok := query[f]; ok {
			v := vv[0]
			switch filter {
			case STATUS_EQ:
				ps = append(ps, jobapplication.StatusEQ(jobapplication.Status(v)))
			case STATUS_NEQ:
				ps = append(ps, jobapplication.StatusNEQ(jobapplication.Status(v)))
			default:
				errors = append(errors, fmt.Errorf("%s is not a valid filter", filter))
			}
		}
	}
	records, err := s.JobApplicationRepository.Filter(ps...)
	if err != nil {
		errors = append(errors, err)
	}
	return records, errors
}
