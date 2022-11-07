package query

import (
	"os"
	"strconv"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

type AlgoliaService struct {
	client    *search.Client
	initIndex *search.Index
}

func NewAlgoliaSearch() *AlgoliaService {
	appID := os.Getenv("ALGOLIA_APP_ID")
	apiKey := os.Getenv("ALGOLIA_API_KEY")
	index := "talent_index"
	client := search.NewClient(appID, apiKey)
	initIndex := client.InitIndex(index)
	return &AlgoliaService{
		client:    client,
		initIndex: initIndex,
	}
}

func (alg *AlgoliaService) Query(params SearchParams) (SearchResult, error) {
	page := 0
	if params.Page != "" {
		p, err := strconv.Atoi(params.Page)
		if err != nil {
			return SearchResult{}, err
		}
		page = p
	}

	combineParams := []interface{}{
		opt.AttributesToRetrieve(params.AttributeToRetrieve...),
		opt.HitsPerPage(params.Limit),
		opt.Page(page),
	}
	res, err := alg.initIndex.Search(params.Text, combineParams...)
	if err != nil {
		return SearchResult{}, err
	}

	result := SearchResult{
		Total:         res.NbPages,
		ItemsThisPage: res.HitsPerPage,
		Items:         res.Hits,
	}

	return result, nil
}
