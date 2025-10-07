package query

import (
	"errors"
	"os"
	"strconv"

	"github.com/theterminalguy/tentn/internal/paginator"
	"github.com/theterminalguy/tentn/internal/repository"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/google/uuid"
)

type AlgoliaService struct {
	client     *search.Client
	initIndex  *search.Index
	talentRepo repository.TalentRepository
}

func NewAlgoliaSearch() *AlgoliaService {
	appID := os.Getenv("ALGOLIA_APP_ID")
	apiKey := os.Getenv("ALGOLIA_API_KEY")
	index := "talent_search_index_schema"
	client := search.NewClient(appID, apiKey)
	initIndex := client.InitIndex(index)
	return &AlgoliaService{
		client:     client,
		initIndex:  initIndex,
		talentRepo: *repository.NewTalentRepository(),
	}
}

func (alg *AlgoliaService) Query(params *SearchParams) (*paginator.OffsetPaginater, []error) {
	page := 0
	var errs []error
	if params.Page != "" {
		p, err := strconv.Atoi(params.Page)
		if err != nil {
			errs = append(errs, err)
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
		errs = append(errs, err)
	}

	if res.NbPages == 0 {
		errs = append(errs, errors.New("no record found"))
	}

	userIds, err := getUserIds(res)
	if err != nil {
		errs = append(errs, err)
	}

	talents, err := alg.talentRepo.GetTalentByUserIDs(userIds...)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	next := 0
	if page < res.NbPages {
		next = page + 1
	} else {
		next = res.NbPages
	}
	result := &paginator.OffsetPaginater{
		Total:         res.NbPages,
		ItemsThisPage: res.HitsPerPage,
		NextCursor:    strconv.Itoa(next),
		Items:         talents,
	}

	return result, nil
}

func getUserIds(records search.QueryRes) ([]interface{}, error) {
	userIds := make([]interface{}, records.HitsPerPage)
	for k, record := range records.Hits {
		if _, ok := record["objectID"]; ok {
			rid := record["objectID"].(string)
			uId := uuid.MustParse(rid)
			userIds[k] = uId
		}
	}

	return userIds, nil
}
