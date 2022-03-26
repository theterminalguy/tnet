package recruiter

import (
	"fmt"
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/search"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type V1RecruiterJobHandler struct {
	JobRepository repo.JobQuerier
	JobSearch     *search.JobSearch
}

func NewV1RecruiterJobHandler(jobQuerier repo.JobQuerier) *V1RecruiterJobHandler {
	return &V1RecruiterJobHandler{
		JobRepository: jobQuerier,
	}
}

func (h *V1RecruiterJobHandler) Search(c echo.Context) error {
	jobSearch := new(search.JobSearch)
	query := c.QueryString()
	records, vldErrs := jobSearch.Search(query)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *V1RecruiterJobHandler) ReadAll(c echo.Context) error {
	user, err := GetCurrentRecruiter(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	jobs, err := user.GetJobs()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, jobs)
}

// ReadByID return a job by its id. The job must be created by the recruiter
func (h *V1RecruiterJobHandler) ReadByID(c echo.Context) error {
	user, err := GetCurrentRecruiter(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	j, err := user.GetJobByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, j)
}

// CreateOne creates a new job for the recruiter
func (h *V1RecruiterJobHandler) CreateOne(c echo.Context) error {
	user, err := GetCurrentRecruiter(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.JobParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	params.UserID = user.Recruiter.ID
	j, err := h.JobRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, j)
}

// UpdateByID updates a job by its id. The job must be created by the recruiter
func (h *V1RecruiterJobHandler) UpdateByID(c echo.Context) error {
	user, err := GetCurrentRecruiter(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.JobParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	job, err := user.GetJobByID(id)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	j, vldErrs := h.JobRepository.Update(job.UUID, *params)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, j)
}

// DeleteByID deletes a job by its id. The job must be created by the recruiter
func (h *V1RecruiterJobHandler) DeleteOne(c echo.Context) error {
	user, err := GetCurrentRecruiter(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	job, err := user.GetJobByID(id)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.JobRepository.DeleteByID(job.UUID)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
