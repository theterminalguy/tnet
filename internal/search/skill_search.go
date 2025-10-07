package search

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/theterminalguy/tnet/ent"
	"github.com/theterminalguy/tnet/ent/predicate"
	"github.com/theterminalguy/tnet/ent/skill"
	repo "github.com/theterminalguy/tnet/internal/repository"
)

type SkillSearch struct {
	SkillRepository repo.SkillRepository
}

func (*SkillSearch) PossibleFilters() []Filter {
	return []Filter{
		NAME_EQ,
	}
}

func (s *SkillSearch) Search(qs string) ([]*ent.Skill, []error) {
	// TODO: this implementation of search and filters is not reusable but works really well for now
	// and makes it easy for us to decide what can be searchable and what can't.
	// We should probably define a Searchable interface and implement it for each searchable entity
	var ps []predicate.Skill
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
				ps = append(ps, skill.ID(uuid.MustParse(v)))
			case UUID_NEQ:
				ps = append(ps, skill.ID(uuid.MustParse(v)))
			case NAME_EQ:
				ps = append(ps, skill.NameEQ(v))
			case NAME_NEQ:
				ps = append(ps, skill.NameNEQ(v))
			case YEAR_EXP_EQ:
				fv, err := stringToFloat32(v)
				if err != nil {
					errors = append(errors, err)
				}
				ps = append(ps, skill.YearsOfExperienceEQ(fv))
			case YEAR_EXP_NEQ:
				fv, err := stringToFloat32(v)
				if err != nil {
					errors = append(errors, err)
				}
				ps = append(ps, skill.YearsOfExperienceNEQ(fv))
			default:
				errors = append(errors, fmt.Errorf("%s is not a valid filter", filter))
			}
		}
	}
	records, err := s.SkillRepository.Filter(ps...)
	if err != nil {
		errors = append(errors, err)
	}
	return records, errors
}

func stringToFloat32(val string) (float32, error) {
	fv, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return 0.0, err
	}
	return float32(fv), nil
}
