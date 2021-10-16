package handler

import (
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type WorkExperienceHandler struct {
	WorkExperienceService    *service.WorkExperienceService
	WorkExperienceRepository *repo.WorkExperienceRepository
}

func NewWorkExperienceHandler() *WorkExperienceHandler {
	return &WorkExperienceHandler{
		WorkExperienceService:    service.NewWorkExperienceService(),
		WorkExperienceRepository: repo.NewWorkExperienceRepository(),
	}
}

func (*WorkExperienceHandler) ResourceName() string {
	return "applicants/work-experiences"
}

func (h *WorkExperienceHandler) ReadAll(c echo.Context) error {
	records, err := h.WorkExperienceRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *WorkExperienceHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.WorkExperienceRepository.GetByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *WorkExperienceHandler) CreateOne(c echo.Context) error {
	params := new(repo.WorkExperienceParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.WorkExperienceRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *WorkExperienceHandler) UpdateByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.WorkExperienceParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	record, err := h.WorkExperienceRepository.Update(id, *params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *WorkExperienceHandler) DeleteOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.WorkExperienceRepository.DeleteByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
