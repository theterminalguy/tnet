package handler

import (
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
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

func (h *ApplicantHandler) ReadAll(c echo.Context) error {
	a, err := h.ApplicantRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, a)
}

func (h *ApplicantHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	a, err := h.ApplicantRepository.GetByUUID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, a)
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
