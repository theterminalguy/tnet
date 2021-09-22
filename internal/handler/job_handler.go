package handler

import (
	"net/http"

	"github.com/10hourlabs/tentn/oneword"
	"github.com/labstack/echo"
)

type JobHandler struct{}

type JobCreateParams struct {
	Hiring       bool     `json:"hiring"`
	Title        string   `json:"title"`
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
	return c.String(http.StatusOK, "GET /ReadAll")
}

func (JobHandler) ReadByID(c echo.Context) error {
	return c.String(http.StatusOK, "GET /ReadByID")
}

func (JobHandler) CreateOne(c echo.Context) error {
	j := new(JobCreateParams)
	if err := c.Bind(j); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, j)
}

func (JobHandler) UpdateByID(c echo.Context) error {
	return c.String(http.StatusOK, "PUT /UpdateByID")
}

func (JobHandler) DeleteOne(c echo.Context) error {
	return c.String(http.StatusNoContent, "DELETE /DeleteOne")
}
