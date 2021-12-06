package recruiter

import (
	"fmt"
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type V1RecruiterJobHandler struct {
	JobRepository *repo.JobRepository
}

func NewV1RecruiterJobHandler() *V1RecruiterJobHandler {
	return &V1RecruiterJobHandler{
		JobRepository: repo.NewJobRepository(),
	}
}

// ReadAll returns all jobs created by the recruiter
func (h *V1RecruiterJobHandler) ReadAll(c echo.Context) error {
	// TODO: implement pagination
	// most likely coursor based
	jobs, err := h.JobRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, jobs)
}

// ReadByID return a job by its id. The job must be created by the recruiter
func (h *V1RecruiterJobHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	j, err := h.JobRepository.GetByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, j)
}

// CreateOne creates a new job for the recruiter
func (h *V1RecruiterJobHandler) CreateOne(c echo.Context) error {
	params := new(repo.JobParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	j, err := h.JobRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, j)
}

// UpdateByID updates a job by its id. The job must be created by the recruiter
func (h *V1RecruiterJobHandler) UpdateByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.JobParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	j, vldErrs := h.JobRepository.Update(id, *params)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, j)
}

// DeleteByID deletes a job by its id. The job must be created by the recruiter
func (h *V1RecruiterJobHandler) DeleteOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.JobRepository.DeleteByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
