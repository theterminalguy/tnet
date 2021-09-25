package handler

import (
	"net/http"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/job"
	"github.com/10hourlabs/tentn/internal/services/job_service"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/labstack/echo"
)

type JobHandler struct{}

type JobCreateParams struct {
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

func (JobHandler) BasePath() string {
	return oneword.Jobs
}

func (JobHandler) ReadAll(c echo.Context) error {
	jobs, err := job_service.GetAllJobs()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, jobs)
}

func (JobHandler) ReadByID(c echo.Context) error {
	return c.String(http.StatusOK, "GET /ReadByID")
}

func (JobHandler) CreateOne(c echo.Context) error {
	params := new(JobCreateParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	j := ent.Job{
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
	job, err := job_service.CreateJob(&j)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, job)
}

func (JobHandler) UpdateByID(c echo.Context) error {
	return c.String(http.StatusOK, "PUT /UpdateByID")
}

func (JobHandler) DeleteOne(c echo.Context) error {
	return c.String(http.StatusNoContent, "DELETE /DeleteOne")
}
