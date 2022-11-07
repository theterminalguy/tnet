package query

import (
	"github.com/10hourlabs/tenlog"
)

type SearchServiceHandler interface {
	Query(params SearchParams) (SearchResult, error)
}

type SearchParams struct {
	Limit               int
	Page                string
	Text                string
	AttributeToRetrieve []string
}

type SearchResult struct {
	Total         int                      `json:"total"`
	ItemsThisPage int                      `json:"items_this_page"`
	Items         []map[string]interface{} `json:"items"`
}

func NewSearchService(query_type string) SearchServiceHandler {
	if query_type == "" {
		tenlog.Error("Unable to process query engine")
		return nil
	}

	query := map[string]SearchServiceHandler{
		"algolia": NewAlgoliaSearch(),
		"rql":     NewRQLService(),
	}

	return query[query_type]
}

func Logger() error {
	return nil
}
