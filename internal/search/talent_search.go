package search

import (
	"fmt"
	"net/url"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/predicate"
	"github.com/10hourlabs/tentn/ent/talent"
	repo "github.com/10hourlabs/tentn/internal/repository"
)

type TalentSearch struct {
	TalentRepository repo.TalentRepository
}

func (*TalentSearch) PossibleFilters() []Filter {
	return []Filter{
		EMAIL_EQ,

		REFRERRAL_CODE_EQ,	
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
			case TENTN_CODE_EQ:
				ps = append(ps, talent.TentnCodeEQ(v))
			case REFRERRAL_CODE_EQ:
				ps = append(ps, talent.ReferralCodeEQ(v))
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
