package talent

import (
	"fmt"
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/search"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type V1JobApplicationHandler struct {
	JobApplicationService    *service.JobApplicationService
	JobApplicationRepository repo.JobApplicationQuerier
}

func NewV1JobApplicationHandler(jobAppQuerier repo.JobApplicationQuerier) *V1JobApplicationHandler {
	return &V1JobApplicationHandler{
		JobApplicationService:    service.NewJobApplicationService(),
		JobApplicationRepository: jobAppQuerier,
	}
}

func (h *V1JobApplicationHandler) Search(c echo.Context) error {
	jobApplicationSearch := new(search.JobApplicationSearch)
	query := c.QueryString()
	records, vldErrs := jobApplicationSearch.Search(query)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *V1JobApplicationHandler) ReadAll(c echo.Context) error {
	talent, err := GetCurrentTalent(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	records, err := talent.GetJobApplications()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *V1JobApplicationHandler) ReadByID(c echo.Context) error {
	talent, err := GetCurrentTalent(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := talent.GetJobApplicationByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1JobApplicationHandler) CreateOne(c echo.Context) error {
	talent, err := GetCurrentTalent(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.JobApplicationParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params.TalentID = talent.Talent.ID
	record, err := h.JobApplicationRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *V1JobApplicationHandler) UpdateByID(c echo.Context) error {
	return nil
}

func (h *V1JobApplicationHandler) DeleteOne(c echo.Context) error {
	return nil
}
