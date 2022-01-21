package search

import (
	"fmt"
	"net/url"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/predicate"
	"github.com/10hourlabs/tentn/ent/skill"
	"github.com/10hourlabs/tentn/ent/talent"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/util/collection"
)

type TalentSearch struct {
	TalentRepository repo.TalentRepository
}

func (*TalentSearch) PossibleFilters() map[string]Filter {
	return map[string]Filter{
		"email":                  EMAIL,
		"email_eq":               EMAIL_EQ,
		"city":                   CITY,
		"city_eq":                CITY_EQ,
		"country":                COUNTRY,
		"years_of_experience":    YEARS_OF_EXPERIENCE,
		"years_of_experience_eq": YEARS_OF_EXPERIENCE_EQ,
		"skills_in":              SKILLS_IN,
		"preferred_title_like":   PREFERRED_TITLE_LIKE,
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
	// loop through all user provided filters
	for f := range query {
		if s.PossibleFilters()[f] == "" {
			errors = append(errors, fmt.Errorf("%s is not a valid filter", f))
			continue
		}
	}
	for _, filter := range s.PossibleFilters() {
		f := string(filter)
		if values, ok := query[f]; ok {
			v := values[0]
			switch filter {
			case CITY, CITY_EQ:
				// TODO: this should be a blog post
				ps = append(ps, talent.CityEqualFold(v))
			case EMAIL, EMAIL_EQ:
				ps = append(ps, talent.EmailEqualFold(v))
			case COUNTRY:
				ps = append(ps, talent.CountryCodeEqualFold(v))
			case YEARS_OF_EXPERIENCE:
				//TODO: This should be a blog post
				filter := func() predicate.Talent {
					return predicate.Talent(func(s *sql.Selector) {
						s.Where(sql.ExprP("DATE_PART('year', AGE(CURRENT_DATE, professional_start_date)) >= $1", v))
					})
				}
				ps = append(ps, filter())
			case YEARS_OF_EXPERIENCE_EQ:
				filter := func() predicate.Talent {
					return predicate.Talent(func(s *sql.Selector) {
						s.Where(sql.ExprP("DATE_PART('year', AGE(CURRENT_DATE, professional_start_date)) = $1", v))
					})
				}
				ps = append(ps, filter())
			case SKILLS_IN:
				// TODO: check if name on skills table is indexed
				sks := strings.Split(v, ",")
				sksPred := skill.NameIn(sks...)
				ps = append(ps, talent.HasSkillsWith(sksPred))
			case PREFERRED_TITLE_LIKE:
				ps = append(ps, talent.PreferredJobTitleContainsFold("%"+v+"%"))
			}
		}
	}
	if collection.HasAny(errors) {
		return nil, errors
	}
	records, err := s.TalentRepository.Filter(ps...)
	if err != nil {
		errors = append(errors, err)
		return nil, errors
	}
	return records, nil
}
