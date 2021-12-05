package search

import (
	"fmt"
	"net/url"

	//"time"

	//"strconv"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/predicate"
	"github.com/10hourlabs/tentn/ent/talent"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/google/uuid"
)

type TalentSearch struct {
	TalentRepository repo.TalentRepository
}

func (*TalentSearch) PossibleFilters() []Filter {
	// Terrible code but it works
	// the compiler does not have to figure our the type at runtime
	return []Filter{
		UUID_EQ,
		UUID_NEQ,

		EMAIL_EQ,

		COUNTRY_CODE_EQ,
		COUNTRY_CODE_NEQ,

		TENTN_CODE_EQ,

		START,
		END,
	}
}

func (s *TalentSearch) Search(qs string) ([]*ent.Talent, []error) {
	// TODO: this implementation of search and filters is not reusable but works really well for now
	// and makes it easy for us to decide what can be searchable and what can't.
	// We should probably define a Searchable interface and implement it for each searchable entity
	var ps []predicate.Talent
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
				ps = append(ps, talent.UUIDEQ(uuid.MustParse(v)))
			case UUID_NEQ:
				ps = append(ps, talent.UUIDNEQ(uuid.MustParse(v)))
			case EMAIL_EQ:
				ps = append(ps, talent.EmailEQ(v))
			case COUNTRY_CODE_EQ:
				ps = append(ps, talent.CountryCodeEQ(v))
			case COUNTRY_CODE_NEQ:
				ps = append(ps, talent.CountryCodeNEQ(v))
			case TENTN_CODE_EQ:
				ps = append(ps, talent.TentnCodeEQ(v))
			// case START:
			// 	t, err := time.Parse("2006-01-02", v)
			// 	if err != nil {
			// 		errors = append(errors, err)
			// 	}
			// 	ps = append(ps, talent.CreatedAtGTE(t))
			// case END:
			// 	t, err := time.Parse("2006-01-02", v)
			// 	if err != nil {
			// 		errors = append(errors, err)
			// 	}
			// 	ps = append(ps, talent.CreatedAtLTE(t))

			default:
				errors = append(errors, fmt.Errorf("%s is not a valid filter", filter))
			}
		}
	}
	records, err := s.TalentRepository.Filter(ps...)
	if err != nil {
		errors = append(errors, err)
	}
	return records, errors
}
