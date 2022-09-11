package recruiter

import (
	"fmt"
	"net/http"
	"os"

	"github.com/10hourlabs/tenlog"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/repository/scope"
	"github.com/10hourlabs/tentn/internal/search"
	"github.com/10hourlabs/tentn/internal/service/payment"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type V1RecruiterJobHandler struct {
	JobRepository        repo.JobQuerier
	JobSearch            *search.JobSearch
	JobRepo              *repo.JobRepository
	TalentCollectionRepo repo.TalentCollectionRepository
}

func NewV1RecruiterJobHandler(jobQuerier repo.JobQuerier) *V1RecruiterJobHandler {
	return &V1RecruiterJobHandler{
		JobRepository:        jobQuerier,
		JobRepo:              repo.NewJobRepository(),
		TalentCollectionRepo: *repo.NewTalentCollectionRepository(),
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
	params := new(repo.JobParams)
	if err := c.Bind(params); err != nil {
		tenlog.Error(err.Error())
		return c.String(http.StatusBadRequest, "Unable to process payload, Please try again later")
	}
	if params.Summary == "" {
		params.Summary = "N/A"
	}
	if params.Thumbnail == "" {
		params.Thumbnail = "https://"
	}
	if params.TimeZone == "" {
		params.TimeZone = "GMT"
	}
	if params.Employment == "" {
		params.Employment = "full_time"
	}
	if params.Category == "" {
		params.Category = "engineering"
	}
	params.UserID = recruiterID

	// Create collection
	jd, err := h.JobRepo.Create(*params)
	if err != nil {
		tenlog.Error(err.Error())
		return c.String(http.StatusBadRequest, "error occured while processing job: "+err.Error())
	}

	// Generate payment link
	driver := os.Getenv("PAYMENT_DRIVER")
	pay := payment.NewPaymentService(driver)
	_, err = pay.GenerateLink(jd.ID)
	if err != nil {
		tenlog.Error(err.Error())
		return c.String(http.StatusBadRequest, "error occured while generate payment link: "+err.Error())
	}

	response, err := currentRecruiter.GetJobByID(jd.ID)
	if err != nil {
		tenlog.Error(err.Error())
		return c.String(http.StatusBadRequest, "error occured while processing JD: "+err.Error())
	}
	return c.JSON(http.StatusCreated, response)
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
