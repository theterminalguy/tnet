package recruiter

import (
	"fmt"
	"net/http"

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
	query := c.QueryString()
	records, vldErrs := talentSearch.Search(query)
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
