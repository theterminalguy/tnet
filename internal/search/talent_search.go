package search

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/10hourlabs/tentn/ent/predicate"
	"github.com/10hourlabs/tentn/ent/skill"
	"github.com/10hourlabs/tentn/ent/talent"
	"github.com/10hourlabs/tentn/internal/paginator"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/util"
	"github.com/10hourlabs/tentn/util/collection"
)

type TalentSearch struct {
	TalentRepository repo.TalentRepository
}

func (*TalentSearch) PossibleFilters() map[Filter]Filter {
	return map[Filter]Filter{
		EMAIL:                    EMAIL,
		EMAIL_EQ:                 EMAIL_EQ,
		CITY:                     CITY,
		CITY_EQ:                  CITY_EQ,
		COUNTRY:                  COUNTRY,
		COUNTRY_EQ:               COUNTRY_EQ,
		SKILLS_EQ:                SKILLS_EQ,
		SKILLS_IN:                SKILLS_IN,
		LOCATED_IN:               LOCATED_IN,
		FIRST_NAME_EQ:            FIRST_NAME_EQ,
		LAST_NAME_EQ:             LAST_NAME_EQ,
		JOB_TITLE_EQ:             JOB_TITLE_EQ,
		YEARS_OF_EXPERIENCE_EQ:   YEARS_OF_EXPERIENCE_EQ,
		YEARS_OF_EXPERIENCE_LT:   YEARS_OF_EXPERIENCE_LT,
		YEARS_OF_EXPERIENCE_GT:   YEARS_OF_EXPERIENCE_GT,
		YEARS_OF_EXPERIENCE_GTEQ: YEARS_OF_EXPERIENCE_GTEQ,
		YEARS_OF_EXPERIENCE_LTEQ: YEARS_OF_EXPERIENCE_LTEQ,
		IS_AVAILABLE:             IS_AVAILABLE,
		JOB_PREFERENCE:           JOB_PREFERENCE,
	}
}

func (*TalentSearch) yearsOfExpereinceFilter(s string, op Operator) (predicate.Talent, error) {
	// TODO: this is should be a blog post
	s = util.ExtractFirstNumbers(s)
	if s == "" {
		return nil, fmt.Errorf("provide a valid years of experience")
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Println("error converting years of experience to float64", err)
		return nil, err
	}
	expr := ""
	switch op {
	case EQ:
		expr = "DATE_PART('year', AGE(CURRENT_DATE, professional_start_date)) = $1"
	case LT:
		expr = "DATE_PART('year', AGE(CURRENT_DATE, professional_start_date)) < $1"
	case GT:
		expr = "DATE_PART('year', AGE(CURRENT_DATE, professional_start_date)) > $1"
	case GTEQ:
		expr = "DATE_PART('year', AGE(CURRENT_DATE, professional_start_date)) >= $1"
	case LTEQ:
		expr = "DATE_PART('year', AGE(CURRENT_DATE, professional_start_date)) <= $1"
	}
	return predicate.Talent(func(s *sql.Selector) {
		s.Where(sql.ExprP(expr, v))
	}), nil
}

func (*TalentSearch) jobPreferenceFilter(v string) (predicate.Talent, error) {
	return predicate.Talent(func(s *sql.Selector) {
		s.Where(sqljson.ValueContains(talent.FieldJobPreference, v))
	}), nil
}

func (s *TalentSearch) Search(cursor, qs string) (*paginator.OffsetPaginater, []error) {
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
		filter := Filter(f)
		if filter == "cursor" {
			continue
		}
		if s.PossibleFilters()[filter] == "" {
			errors = append(errors, fmt.Errorf("%s is not a valid filter", f))
			continue
		}
	}
	for _, possibleFilter := range s.PossibleFilters() {
		queryFilter := string(possibleFilter)
		if values, ok := query[queryFilter]; ok {
			v := values[0]
			switch possibleFilter {
			case CITY, CITY_EQ:
				// TODO: this should be a blog post
				ps = append(ps, talent.CityEqualFold(v))
			case EMAIL, EMAIL_EQ:
				ps = append(ps, talent.EmailEqualFold(v))
			case COUNTRY, COUNTRY_EQ:
				ps = append(ps, talent.CountryCodeEqualFold(v))
			case YEARS_OF_EXPERIENCE_GTEQ:
				filter, err := s.yearsOfExpereinceFilter(v, GTEQ)
				if err != nil {
					errors = append(errors, err)
					continue
				}
				ps = append(ps, filter)
			case YEARS_OF_EXPERIENCE_EQ:
				filter, err := s.yearsOfExpereinceFilter(v, EQ)
				if err != nil {
					errors = append(errors, err)
					continue
				}
				ps = append(ps, filter)
			case YEARS_OF_EXPERIENCE_LT:
				filter, err := s.yearsOfExpereinceFilter(v, LT)
				if err != nil {
					errors = append(errors, err)
					continue
				}
				ps = append(ps, filter)
			case YEARS_OF_EXPERIENCE_LTEQ:
				filter, err := s.yearsOfExpereinceFilter(v, LTEQ)
				if err != nil {
					errors = append(errors, err)
					continue
				}
				ps = append(ps, filter)
			case SKILLS_EQ:
				ps = append(ps, talent.HasSkillsWith(skill.NameEqualFold(v)))
			case SKILLS_IN:
				// TODO: check if name on skills table is indexed
				sks := strings.Split(strings.ToLower(v), ",")
				sksPred := skill.NameIn(sks...)
				ps = append(ps, talent.HasSkillsWith(sksPred))
			case JOB_TITLE_EQ:
				ps = append(ps, talent.PreferredJobTitleContainsFold("%"+v+"%"))
			case FIRST_NAME_EQ:
				ps = append(ps, talent.FirstNameEqualFold(v))
			case LAST_NAME_EQ:
				ps = append(ps, talent.LastNameEqualFold(v))
			case LOCATED_IN:
				var city, country string
				values := strings.Split(v, ",")
				if len(values) == 2 {
					city = values[0]
					country = values[1]
					cityPred := talent.CityEqualFold(city)
					countryPred := talent.CountryCodeEqualFold(country)
					ps = append(ps, cityPred, countryPred)
				}
			case IS_AVAILABLE:
				booleanField, err := strconv.ParseBool(v)
				if err != nil {
					errors = append(errors, err)
				}
				ps = append(ps, talent.IsAvailableEQ(booleanField))
			case JOB_PREFERENCE:
				filter, err := s.jobPreferenceFilter(v)
				if err != nil {
					errors = append(errors, err)
					continue
				}
				ps = append(ps, filter)
			}
		}
	}
	if collection.HasAny(errors) {
		return nil, errors
	}
	ps = append(ps, talent.HasEducations(), talent.HasSkills(), talent.HasWorkExperiences(), talent.HasPortfoliolinks())
	records, err := s.TalentRepository.Filter(cursor, ps...)
	if err != nil {
		errors = append(errors, err)
		return nil, errors
	}
	return records, nil
}
