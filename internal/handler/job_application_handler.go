package handler

import (
	"fmt"
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type JobApplicationHandler struct {
	JobApplicationService    *service.JobApplicationService
	JobApplicationRepository *repo.JobApplicationRepository
}

func NewJobApplicationHandler() *JobApplicationHandler {
	return &JobApplicationHandler{
		JobApplicationService:    service.NewJobApplicationService(),
		JobApplicationRepository: repo.NewJobApplicationRepository(),
	}
}

func (*JobApplicationHandler) ResourceName() string {
	return "jobs/applications"
}

func (h *JobApplicationHandler) ReadAll(c echo.Context) error {
	records, err := h.JobApplicationRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *JobApplicationHandler) ReadByID(c echo.Context) error {
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

func (h *JobApplicationHandler) CreateOne(c echo.Context) error {
	params := new(repo.JobApplicationParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.JobApplicationRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *JobApplicationHandler) UpdateByID(c echo.Context) error {
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

func (h *JobApplicationHandler) DeleteOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.JobApplicationRepository.DeleteByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
