package search

import (
	"fmt"
	"net/url"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/education"
	"github.com/10hourlabs/tentn/ent/predicate"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/google/uuid"
)

type EducationSearch struct {
	EducationRepository repo.EducationRepository
}

func (*EducationSearch) PossibleFilters() []Filter {
	// Terrible code but it works
	// the compiler does not have to figure our the type at runtime
	return []Filter{
		UUID_EQ,
		UUID_NEQ,

		DEGREE_EQ,
		DEGREE_NEQ,

		INST_NAME_EQ,
		INST_NAME_NEQ,

		PROGRAM_EQ,
		PROGRAM_NEQ,
	}
}

func (s *EducationSearch) Search(qs string) ([]*ent.Education, []error) {
	// TODO: this implementation of search and filters is not reusable but works really well for now
	// and makes it easy for us to decide what can be searchable and what can't.
	// We should probably define a Searchable interface and implement it for each searchable entity
	var ps []predicate.Education
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
				ps = append(ps, education.UUIDEQ(uuid.MustParse(v)))
			case UUID_NEQ:
				ps = append(ps, education.UUIDNEQ(uuid.MustParse(v)))
			case DEGREE_EQ:
				ps = append(ps, education.DegreeEQ(v))
			case DEGREE_NEQ:
				ps = append(ps, education.DegreeNEQ(v))
			case INST_NAME_EQ:
				ps = append(ps, education.InstitutionNameEQ(v))
			case INST_NAME_NEQ:
				ps = append(ps, education.InstitutionNameNEQ(v))
			case PROGRAM_EQ:
				ps = append(ps, education.ProgramEQ(v))
			case PROGRAM_NEQ:
				ps = append(ps, education.ProgramNEQ(v))
			default:
				errors = append(errors, fmt.Errorf("%s is not a valid filter", filter))
			}
		}
	}
	records, err := s.EducationRepository.Filter(ps...)
	if err != nil {
		errors = append(errors, err)
	}
	return records, errors
}
