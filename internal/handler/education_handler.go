package handler

import (
	"fmt"
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/search"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type EducationHandler struct {
	EducationService    *service.EducationService
	EducationRepository *repo.EducationRepository
}

func NewEducationHandler() *EducationHandler {
	return &EducationHandler{
		EducationService:    service.NewEducationService(),
		EducationRepository: repo.NewEducationRepository(),
	}
}

func (h *EducationHandler) Search(c echo.Context) error {
	educationSearch := new(search.EducationSearch)
	query := c.QueryString()
	records, vldErrs := educationSearch.Search(query)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *EducationHandler) ReadAll(c echo.Context) error {
	records, err := h.EducationRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *EducationHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.EducationRepository.GetByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *EducationHandler) CreateOne(c echo.Context) error {
	params := new(repo.EducationParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.EducationRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *EducationHandler) UpdateByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.EducationParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	record, vldErrs := h.EducationRepository.Update(id, *params)
	if vldErrs != nil {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *EducationHandler) DeleteOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.EducationRepository.DeleteByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
