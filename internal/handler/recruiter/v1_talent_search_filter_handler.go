package recruiter

import (
	"fmt"
	"net/http"

	"github.com/10hourlabs/rql/parser"
	"github.com/10hourlabs/tentn/ent"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/search"
	"github.com/labstack/echo/v4"
)

type V1TalentSearchFilterHandler struct {
	TalentRepo *repo.TalentRepository
}

func NewV1TalentSearchFilterHandler() *V1TalentSearchFilterHandler {
	return &V1TalentSearchFilterHandler{
		TalentRepo: repo.NewTalentRepository(),
	}
}

func (h *V1TalentSearchFilterHandler) Search(c echo.Context) error {
	talentSearch := new(search.TalentSearch)
	fmt.Println(c.QueryParams())
	query := c.QueryString()
	if query == "" {
		// TODO: I am not sure if this should be the correct thing to do
		return c.JSON(http.StatusBadRequest, "Provide at least one filter")
	}
	var records []*ent.Talent
	var vldErrs []error

	if c.QueryParams().Has("q") {
		// use the RQL parser here
		query, err := parser.Eval(c.QueryParam("q"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		records, vldErrs = talentSearch.Search(query)
	} else {
		records, vldErrs = talentSearch.Search(query)
	}
	if vldErrs != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": fmt.Errorf("%v", vldErrs).Error()})
	}
	return c.JSON(http.StatusOK, records)
}

func (*V1TalentSearchFilterHandler) ReadAll(c echo.Context) error {
	return nil
}

func (*V1TalentSearchFilterHandler) ReadByID(c echo.Context) error {
	return nil
}

func (*V1TalentSearchFilterHandler) CreateOne(c echo.Context) error {
	return nil
}

func (*V1TalentSearchFilterHandler) UpdateByID(c echo.Context) error {
	return nil
}

func (*V1TalentSearchFilterHandler) DeleteOne(c echo.Context) error {
	return nil
}
