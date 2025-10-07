package handler

import (
	"fmt"
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type EmergencyContactHandler struct {
	EmergencyContactService    *service.EmergencyContactService
	EmergencyContactRepository *repo.EmergencyContactRepository
}

func NewEmergencyContactHandler() *EmergencyContactHandler {
	return &EmergencyContactHandler{
		EmergencyContactService:    service.NewEmergencyContactService(),
		EmergencyContactRepository: repo.NewEmergencyContactRepository(),
	}
}

func (h *EmergencyContactHandler) ReadAll(c echo.Context) error {
	records, err := h.EmergencyContactRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *EmergencyContactHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.EmergencyContactRepository.GetByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *EmergencyContactHandler) CreateOne(c echo.Context) error {
	params := new(repo.EmergencyContactParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.EmergencyContactRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *EmergencyContactHandler) UpdateByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.EmergencyContactParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	record, vldErrs := h.EmergencyContactRepository.Update(id, *params)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *EmergencyContactHandler) DeleteOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.EmergencyContactRepository.DeleteByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
