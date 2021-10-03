package handler

import (
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/labstack/echo"
)

type ApplicantHandler struct {
	ApplicantService    *service.ApplicantService
	ApplicantRepository *repo.ApplicantRepository
}

func NewApplicantHandler() *ApplicantHandler {
	return &ApplicantHandler{
		ApplicantService:    service.NewApplicantService(),
		ApplicantRepository: repo.NewApplicantRepository(),
	}
}

func (*ApplicantHandler) ResourceName() string {
	return oneword.Applicants
}

func (*ApplicantHandler) ReadAll(c echo.Context) error {
	return c.String(http.StatusOK, "GET /applicants")
}

func (*ApplicantHandler) ReadByID(c echo.Context) error {
	return c.String(http.StatusOK, "GET /applicants/:some-id")
}

func (h *ApplicantHandler) CreateOne(c echo.Context) error {
	params := new(repo.ApplicantParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	a, err := h.ApplicantRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, a)
}

func (*ApplicantHandler) UpdateByID(c echo.Context) error {
	return c.String(http.StatusOK, "PUT /applicants/:some-id")
}

func (*ApplicantHandler) DeleteOne(c echo.Context) error {
	return c.String(http.StatusOK, "DELETE /applicants/:some-id")
}
