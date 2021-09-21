package handler

import (
	"net/http"

	"github.com/10hourlabs/tentn/oneword"
	"github.com/labstack/echo"
)

type JobHandler struct {
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
	return c.String(http.StatusCreated, "POST /CreateOne")
}

func (JobHandler) UpdateByID(c echo.Context) error {
	return c.String(http.StatusOK, "PUT /UpdateByID")
}

func (JobHandler) DeleteOne(c echo.Context) error {
	return c.String(http.StatusNoContent, "DELETE /DeleteOne")
}
