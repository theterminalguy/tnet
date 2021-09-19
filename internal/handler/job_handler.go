package handler

import (
	"github.com/10hourlabs/tentn/oneword"
	"github.com/labstack/echo"
)

type JobHandler struct {
}

func (JobHandler) BasePath() string {
	return oneword.Jobs
}

func (JobHandler) ReadByID(c echo.Context) error {
	return nil
}

func (JobHandler) ReadAll(c echo.Context) error {
	return nil
}

func (JobHandler) CreateOne(c echo.Context) error {
	return nil
}

func (JobHandler) UpdateByID(c echo.Context) error {
	return nil
}

func (JobHandler) DeleteOne(c echo.Context) error {
	return nil
}
