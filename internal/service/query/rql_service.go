package query

type RQLService struct{}

func NewRQLService() *RQLService {
	return &RQLService{}
}

func (*RQLService) Query(params SearchParams) (SearchResult, error) {
	// talentSearch := new(search.TalentSearch)
	// records, _ = talentSearch.Search(params.Page, params.Text)
	// if vldErrs != nil {
	// 	return nil, vldErrs
	// }
	return SearchResult{}, nil
}
