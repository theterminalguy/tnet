package recruiter

import (
	"fmt"
	"net/http"

	"github.com/10hourlabs/rql/parser"
	"github.com/10hourlabs/tentn/internal/paginator"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/search"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
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
	query := c.QueryString()
	var records *paginator.OffsetPaginater
	var vldErrs []error

	page := c.QueryParam("cursor")
	if c.QueryParams().Has("q") && c.QueryParams().Get("q") != "" {
		// use the RQL parser here
		query, err := parser.Eval(c.QueryParam("q"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		records, vldErrs = talentSearch.Search(page, query)
	} else {
		records, vldErrs = talentSearch.Search(page, query)
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

func (v *V1TalentSearchFilterHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := v.TalentRepo.GetByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
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
