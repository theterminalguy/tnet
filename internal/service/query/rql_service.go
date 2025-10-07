package query

import (
	"github.com/theterminalguy/rql/parser"
	"github.com/theterminalguy/tentn/internal/paginator"
	"github.com/theterminalguy/tentn/internal/search"
)

type RQLService struct{}

func NewRQLService() *RQLService {
	return &RQLService{}
}

func (*RQLService) Query(params *SearchParams) (*paginator.OffsetPaginater, []error) {
	var errors []error
	query, err := parser.Eval(params.Text)
	if err != nil {
		errors = append(errors, err)
	}
	talentSearch := new(search.TalentSearch)
	records, vldErrs := talentSearch.Search(params.Page, query)
	if vldErrs != nil {
		errors = append(errors, vldErrs...)
	}
	if len(errors) > 0 {
		return nil, errors
	}
	return records, nil
}
