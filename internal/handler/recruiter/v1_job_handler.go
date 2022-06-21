package recruiter

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/repository/scope"
	"github.com/10hourlabs/tentn/internal/search"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/10hourlabs/tentn/util"
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
	currentRecruiter := c.Get(oneword.CurrentRecruiter).(*scope.RecruiterScope)
	jobs, err := currentRecruiter.GetJobs()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, jobs)
}

// ReadByID return a job by its id. The job must be created by the recruiter
func (h *V1RecruiterJobHandler) ReadByID(c echo.Context) error {
	currentRecruiter := c.Get(oneword.CurrentRecruiter).(*scope.RecruiterScope)
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	j, err := currentRecruiter.GetJobByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, j)
}

// CreateOne creates a new job for the recruiter
func (h *V1RecruiterJobHandler) CreateOne(c echo.Context) error {
	currentRecruiter := c.Get(oneword.CurrentRecruiter).(*scope.RecruiterScope)
	recruiterID := currentRecruiter.GetID()
	const MAX_FILE_SIZE = 1024 * 1024 * 10 // 10MB
	directory := fmt.Sprintf("data/jd/%s", recruiterID)
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return err
	}
	response, err := util.FileUpload(c, directory, MAX_FILE_SIZE)
	if err != nil {
		return err
	}
	if !strings.Contains(response, directory) {
		return c.String(http.StatusOK, response)
	}

	return c.String(http.StatusOK, "Document received")

	// currentRecruiter := c.Get(oneword.CurrentRecruiter).(*scope.RecruiterScope)
	// params := new(repo.JobParams)
	// if err := c.Bind(params); err != nil {
	// 	return err
	// }

	// params.UserID = currentRecruiter.GetID()
	// j, err := h.JobRepository.Create(*params)
	// if err != nil {
	// 	return c.String(http.StatusBadRequest, err.Error())
	// }
	// return c.JSON(http.StatusCreated, j)
}

// UpdateByID updates a job by its id. The job must be created by the recruiter
func (h *V1RecruiterJobHandler) UpdateByID(c echo.Context) error {
	currentRecruiter := c.Get(oneword.CurrentRecruiter).(*scope.RecruiterScope)
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.JobParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	job, err := currentRecruiter.GetJobByID(id)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	j, vldErrs := h.JobRepository.Update(job.ID, *params)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, j)
}

// DeleteByID deletes a job by its id. The job must be created by the recruiter
func (h *V1RecruiterJobHandler) DeleteOne(c echo.Context) error {
	currentRecruiter := c.Get(oneword.CurrentRecruiter).(*scope.RecruiterScope)
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	job, err := currentRecruiter.GetJobByID(id)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.JobRepository.DeleteByID(job.ID)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
