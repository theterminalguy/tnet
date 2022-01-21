package recruiter

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

type V1RecruiterJobApplicationHandler struct {
	JobApplicationService    *service.JobApplicationService
	JobApplicationRepository repo.JobApplicationQuerier
}

func NewV1RecruiterJobApplicationHandler(jobAppQuerier repo.JobApplicationQuerier) *V1RecruiterJobApplicationHandler {
	return &V1RecruiterJobApplicationHandler{
		JobApplicationService:    service.NewJobApplicationService(),
		JobApplicationRepository: jobAppQuerier,
	}
}

func (h *V1RecruiterJobApplicationHandler) Search(c echo.Context) error {
	jobApplicationSearch := new(search.JobApplicationSearch)
	query := c.QueryString()
	records, vldErrs := jobApplicationSearch.Search(query)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *V1RecruiterJobApplicationHandler) ReadAll(c echo.Context) error {
	// records, err := h.JobApplicationRepository.GetAllForTalent(talent.ID)
	// if err != nil {
	// 	return c.String(http.StatusBadRequest, err.Error())
	// }
	return c.JSON(http.StatusOK, nil)
}

func (h *V1RecruiterJobApplicationHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.JobApplicationRepository.GetByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1RecruiterJobApplicationHandler) CreateOne(c echo.Context) error {
	return nil
}

func (h *V1RecruiterJobApplicationHandler) UpdateByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.JobApplicationParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	record, vldErrs := h.JobApplicationRepository.Update(id, *params)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1RecruiterJobApplicationHandler) DeleteOne(c echo.Context) error {
	return nil
}
