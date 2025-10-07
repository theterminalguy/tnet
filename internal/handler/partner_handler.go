package handler

import (
	"fmt"
	"net/http"

	repo "github.com/theterminalguy/tentn/internal/repository"
	"github.com/theterminalguy/tentn/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PartnerHandler struct {
	PartnerService    *service.PartnerService
	PartnerRepository *repo.PartnerRepository
}

func NewPartnerHandler() *PartnerHandler {
	return &PartnerHandler{
		PartnerService:    service.NewPartnerService(),
		PartnerRepository: repo.NewPartnerRepository(),
	}
}

func (h *PartnerHandler) ReadAll(c echo.Context) error {
	records, err := h.PartnerRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *PartnerHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.PartnerRepository.GetByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *PartnerHandler) CreateOne(c echo.Context) error {
	params := new(repo.PartnerParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.PartnerRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *PartnerHandler) UpdateByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.PartnerParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	record, vldErrs := h.PartnerRepository.Update(id, *params)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *PartnerHandler) DeleteOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.PartnerRepository.DeleteByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
