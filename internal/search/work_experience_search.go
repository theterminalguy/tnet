package search

import (
	"fmt"
	"net/url"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/predicate"
	"github.com/10hourlabs/tentn/ent/workexperience"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/google/uuid"
)

type WorkExperienceSearch struct {
	WorkExperienceRepository repo.WorkExperienceRepository
}

func (*WorkExperienceSearch) PossibleFilters() []Filter {
	// Terrible code but it works
	// the compiler does not have to figure our the type at runtime
	return []Filter{
		UUID_EQ,
		UUID_NEQ,

		COMPANY_NAME_EQ,
		COMPANY_NAME_NEQ,

		LOCATION_EQ,
		LOCATION_NEQ,

		JOB_TITLE_EQ,
		JOB_TITLE_NEQ,

		PRIMARY_TECH_CONT,

		START,
		END,
	}
}

func (s *WorkExperienceSearch) Search(qs string) ([]*ent.WorkExperience, []error) {
	// TODO: this implementation of search and filters is not reusable but works really well for now
	// and makes it easy for us to decide what can be searchable and what can't.
	// We should probably define a Searchable interface and implement it for each searchable entity
	var ps []predicate.WorkExperience
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
				ps = append(ps, workexperience.UUIDEQ(uuid.MustParse(v)))
			case UUID_NEQ:
				ps = append(ps, workexperience.UUIDNEQ(uuid.MustParse(v)))
			case COMPANY_NAME_EQ:
				ps = append(ps, workexperience.CompanyNameEQ(v))
			case COMPANY_NAME_NEQ:
				ps = append(ps, workexperience.CompanyNameNEQ(v))
			case LOCATION_EQ:
				ps = append(ps, workexperience.LocationEQ(v))
			case LOCATION_NEQ:
				ps = append(ps, workexperience.LocationNEQ(v))
			case JOB_TITLE_EQ:
				ps = append(ps, workexperience.JobTitleEQ(v))
			case JOB_TITLE_NEQ:
				ps = append(ps, workexperience.JobTitleNEQ(v))
			case PRIMARY_TECH_CONT:

				//select primary_technologies from work_experiences where primary_technologies @> '["golang"]'

			default:
				errors = append(errors, fmt.Errorf("%s is not a valid filter", filter))
			}
		}
	}
	records, err := s.WorkExperienceRepository.Filter(ps...)
	if err != nil {
		errors = append(errors, err)
	}
	return records, errors
}
