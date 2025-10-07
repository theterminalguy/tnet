package query

import (
	"github.com/theterminalguy/tenlog"
	"github.com/theterminalguy/tentn/internal/paginator"
)

type SearchServiceHandler interface {
	Query(params *SearchParams) (*paginator.OffsetPaginater, []error)
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

var GetSearchInstance SearchServiceHandler

func init() {
	if GetSearchInstance == nil {
		GetSearchInstance = NewSearchService("algolia")
	}
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
