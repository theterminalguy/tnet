package handler

import (
	"fmt"
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type TalentHandler struct {
	TalentService    *service.TalentService
	TalentRepository *repo.TalentRepository
}

func NewTalentHandler() *TalentHandler {
	return &TalentHandler{
		TalentService:    service.NewTalentService(),
		TalentRepository: repo.NewTalentRepository(),
	}
}

func (*TalentHandler) ResourceName() string {
	return oneword.Talents
}

func (h *TalentHandler) ReadAll(c echo.Context) error {
	a, err := h.TalentRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, a)
}

func (h *TalentHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	a, err := h.TalentRepository.GetByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, a)
}

func (h *TalentHandler) CreateOne(c echo.Context) error {
	params := new(repo.TalentParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	a, err := h.TalentRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, a)
}

func (h *TalentHandler) UpdateByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.TalentParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	a, vldErrs := h.TalentRepository.Update(id, *params)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, a)
}

func (h *TalentHandler) DeleteOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.TalentRepository.DeleteByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
