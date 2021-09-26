package handler

import (
	"net/http"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/job"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/labstack/echo"
)

type JobHandler struct {
	JobService *service.JobService
}

func NewJobHandler() *JobHandler {
	// TODO: Decide how to handle failure caused by
	// initializing new service
	// we currently ignore errors caused by dependencies
	// if initialiliz the job service failed
	js, err := service.NewJobService()
	if err != nil {
		// TODO: alternatively, you could add a
		// serviceError field to the JobHandler struct
		// this should get set to true if there was an error initialize a service
		// and you should return internal service error
		return &JobHandler{}
	}
	return &JobHandler{
		JobService: js,
	}
}

type jobCreateParams struct {
	Hiring       bool     `json:"hiring"`
	Title        string   `json:"title"`
	Slug         string   `json:"slug"`
	Summary      string   `json:"summary"`
	Employment   string   `json:"employment"`
	Category     string   `json:"category"`
	Thumbnail    string   `json:"thumbnail"`
	WeHave       []string `json:"we_have"`
	Requirements []string `json:"requirements"`
	YouHave      []string `json:"you_have"`
}

func (*JobHandler) ResourceName() string {
	return oneword.Jobs
}

func (h *JobHandler) ReadAll(c echo.Context) error {
	jobs, err := h.JobService.GetAllJobs()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, jobs)
}

func (*JobHandler) ReadByID(c echo.Context) error {
	return c.String(http.StatusOK, "GET /ReadByID")
}

func (h *JobHandler) CreateOne(c echo.Context) error {
	params := new(jobCreateParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	j := &ent.Job{
		Hiring:       params.Hiring,
		Title:        params.Title,
		Summary:      params.Summary,
		Employment:   job.Employment(params.Employment),
		Category:     job.Category(params.Category),
		Thumbnail:    params.Thumbnail,
		WeHave:       params.WeHave,
		Requirements: params.Requirements,
		YouHave:      params.YouHave,
	}
	j, err := h.JobService.CreateJob(j)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, j)
}

func (*JobHandler) UpdateByID(c echo.Context) error {
	return c.String(http.StatusOK, "PUT /UpdateByID")
}

func (*JobHandler) DeleteOne(c echo.Context) error {
	return c.String(http.StatusNoContent, "DELETE /DeleteOne")
}
