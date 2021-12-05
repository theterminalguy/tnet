package search

import (
	"fmt"
	"net/url"

	//"strconv"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/portfoliolink"
	"github.com/10hourlabs/tentn/ent/predicate"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/google/uuid"
)

type PortfolioLinkSearch struct {
	PortfolioLinkRepository repo.PortfolioLinkRepository
}

func (*PortfolioLinkSearch) PossibleFilters() []Filter {
	// Terrible code but it works
	// the compiler does not have to figure our the type at runtime
	return []Filter{
		UUID_EQ,
		UUID_NEQ,

		NAME_EQ,
		NAME_NEQ,
	}
}

func (s *PortfolioLinkSearch) Search(qs string) ([]*ent.PortfolioLink, []error) {
	// TODO: this implementation of search and filters is not reusable but works really well for now
	// and makes it easy for us to decide what can be searchable and what can't.
	// We should probably define a Searchable interface and implement it for each searchable entity
	var ps []predicate.PortfolioLink
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
				ps = append(ps, portfoliolink.UUIDEQ(uuid.MustParse(v)))
			case UUID_NEQ:
				ps = append(ps, portfoliolink.UUIDNEQ(uuid.MustParse(v)))
			case NAME_EQ:
				ps = append(ps, portfoliolink.NameEQ(v))
			case NAME_NEQ:
				ps = append(ps, portfoliolink.NameNEQ(v))
			default:
				errors = append(errors, fmt.Errorf("%s is not a valid filter", filter))
			}
		}
	}
	records, err := s.PortfolioLinkRepository.Filter(ps...)
	if err != nil {
		errors = append(errors, err)
	}
	return records, errors
}
