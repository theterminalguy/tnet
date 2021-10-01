package handler

import (
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
	return nil
}

func (*ApplicantHandler) ReadByID(c echo.Context) error  {
	return nil
}

func (*ApplicantHandler) CreateOne(c echo.Context) error  {
	return nil
}

func (*ApplicantHandler) UpdateByID(c echo.Context) error  {
	return nil
}

func (*ApplicantHandler) DeleteOne(c echo.Context) error  {
	return nil
}