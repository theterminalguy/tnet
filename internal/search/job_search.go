package search

import (
	"fmt"
	"net/url"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/job"
	"github.com/10hourlabs/tentn/ent/predicate"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/google/uuid"
)

type JobSearch struct {
	JobRepository repo.JobRepository
}

func (*JobSearch) PossibleFilters() []Filter {
	// Terrible code but it works
	// the compiler does not have to figure our the type at runtime
	return []Filter{
		UUID_EQ,
		UUID_NE,

		SLUG_EQ,
		SLUG_NE,

		TITLE_EQ,
		TITLE_NE,
	}
}

func (s *JobSearch) Search(qs string) ([]*ent.Job, []error) {
	// TODO: this implementation of search and filters is not reusable but works really well for now
	// and makes it easy for us to decide what can be searchable and what can't.
	// We should probably define a Searchable interface and implement it for each searchable entity
	var ps []predicate.Job
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
			case UUID_EQ:
				ps = append(ps, job.UUIDEQ(uuid.MustParse(v)))
			case UUID_NE:
				ps = append(ps, job.UUIDNEQ(uuid.MustParse(v)))
			case SLUG_EQ:
				ps = append(ps, job.SlugEQ(v))
			case SLUG_NE:
				ps = append(ps, job.SlugNEQ(v))
			case TITLE_EQ:
				ps = append(ps, job.TitleEQ(v))
			case TITLE_NE:
				ps = append(ps, job.TitleNEQ(v))
			default:
				errors = append(errors, fmt.Errorf("%s is not a valid filter", filter))
			}
		}
	}
	records, err := s.JobRepository.Filter(ps...)
	if err != nil {
		errors = append(errors, err)
	}
	return records, errors
}
