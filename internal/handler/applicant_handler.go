package handler

import (
	"net/http"

	"github.com/10hourlabs/tentn/oneword"
	"github.com/labstack/echo"
)

type ApplicantHandler struct {

}

func NewApplicantHandler() *ApplicantHandler  {
	return &ApplicantHandler{}
}

func (*ApplicantHandler) ResourceName() string {
	return oneword.Applicants
}

func (*ApplicantHandler) ReadAll(c echo.Context) error  {
	return c.String(http.StatusOK, "GET /applicants")
}

func (*ApplicantHandler) ReadByID(c echo.Context) error  {
	return c.String(http.StatusOK, "GET /applicants/:some-id")
}

func (*ApplicantHandler) CreateOne(c echo.Context) error  {
	return c.String(http.StatusCreated, "GET /applicants")
}

func (*ApplicantHandler) UpdateByID(c echo.Context) error  {
	return c.String(http.StatusOK, "PUT /applicants/:some-id")
}

func (*ApplicantHandler) DeleteOne(c echo.Context) error  {
	return c.String(http.StatusOK, "DELETE /applicants/:some-id")
}