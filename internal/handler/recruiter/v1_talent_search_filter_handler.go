package recruiter

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/10hourlabs/tentn/internal/middleware/globalctx"
	"github.com/10hourlabs/tentn/internal/middleware/header"
	"github.com/10hourlabs/tentn/internal/paginator"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/search"
	"github.com/10hourlabs/tentn/internal/service"
	q "github.com/10hourlabs/tentn/internal/service/query"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type V1TalentSearchFilterHandler struct {
	TalentRepo    *repo.TalentRepository
	SearchLogRepo *repo.SearchLogRepository
	PDFService    *service.PDFService
}

func NewV1TalentSearchFilterHandler() *V1TalentSearchFilterHandler {
	return &V1TalentSearchFilterHandler{
		TalentRepo:    repo.NewTalentRepository(),
		SearchLogRepo: repo.NewSearchLogRepository(),
		PDFService:    service.NewPDFService(),
	}
}

func (h *V1TalentSearchFilterHandler) Search(c echo.Context) error {
	query := c.QueryString()
	var records *paginator.OffsetPaginater
	var vldErrs []error

	page := c.QueryParam("cursor")
	if c.QueryParams().Has("q") && c.QueryParams().Get("q") != "" {
		// use the Algolia | RQL parser here
		platform := c.Request().Header.Get(header.X_TN_PLATFORM)
		pageLimit := 3
		if platform == "platform/web" {
			pageLimit = 10
		}
		records, vldErrs = q.GetSearchInstance.Query(&q.SearchParams{
			Limit: pageLimit,
			Page:  page,
			Text:  c.QueryParams().Get("q"),
		})
	} else {
		talentSearch := new(search.TalentSearch)
		records, vldErrs = talentSearch.Search(page, query)
	}
	if vldErrs != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": fmt.Errorf("%v", vldErrs).Error()})
	}
	// Log the search
	searchLog := repo.SearchLogParams{
		Query:          c.QueryParam("q"),
		ResultCount:    records.Total,
		Platform:       globalctx.GetPlatform(c),
		PlatformUserID: globalctx.GetPlatformUserID(c),
		PlatformTeamID: globalctx.GetPlatformTeamID(c),
	}
	if _, err := h.SearchLogRepo.Create(searchLog); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
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
	format := c.QueryParam("format")
	if format == "pdf" {
		c.Response().Header().Set(
			echo.HeaderContentDisposition,
			fmt.Sprintf("attachment; filename=%s-%s.pdf", record.FirstName, record.LastName),
		)
		pdf, err := v.PDFService.Generate(record)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.Blob(http.StatusOK, "application/pdf", pdf)
	}
	if format == "base64" {
		pdf, err := v.PDFService.Generate(record)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]string{"data": base64.StdEncoding.EncodeToString(pdf)})
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
