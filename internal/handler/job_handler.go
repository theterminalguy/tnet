package handler

import (
	"net/http"

	"github.com/10hourlabs/tentn/internal/repo"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type JobHandler struct {
	JobService *service.JobService
	JobRepo    *repo.JobRepository
}

func NewJobHandler() *JobHandler {
	return &JobHandler{
		JobService: service.NewJobService(),
		JobRepo:    repo.NewJobRepository(),
	}
}

func (*JobHandler) ResourceName() string {
	return oneword.Jobs
}

func (h *JobHandler) ReadAll(c echo.Context) error {
	// TODO: implement pagination
	// most likely coursor based
	// also, jobs with hiring = false should NOT
	// be returned
	jobs, err := h.JobRepo.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, jobs)
}

func (h *JobHandler) ReadByID(c echo.Context) error {
	jobUUID := uuid.MustParse(c.Param(oneword.UUID))
	job, err := h.JobRepo.GetByUUID(jobUUID)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, job)
}

func (h *JobHandler) CreateOne(c echo.Context) error {
	params := new(repo.JobParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	j, err := h.JobRepo.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, j)
}

func (h *JobHandler) UpdateByID(c echo.Context) error {
	// TODO: should you be able to update a deleted job?
	jobUUID := uuid.MustParse(c.Param(oneword.UUID))
	params := new(repo.JobParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	job, err := h.JobRepo.Update(jobUUID, *params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, job)
}

func (h *JobHandler) DeleteOne(c echo.Context) error {
	jobUUID, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.JobRepo.DeleteByUUID(jobUUID)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.String(http.StatusNoContent, "")
}
