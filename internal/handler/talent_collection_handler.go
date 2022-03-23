package handler

import (
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TalentCollectionHandler struct {
	TalentCollectionService    *service.TalentCollectionService
	TalentCollectionRepository *repo.TalentCollectionRepository
}

func NewTalentCollectionHandler() *TalentCollectionHandler {
	return &TalentCollectionHandler{
		TalentCollectionService:    service.NewTalentCollectionService(),
		TalentCollectionRepository: repo.NewTalentCollectionRepository(),
	}
}

func (*TalentCollectionHandler) ResourceName() string {
	return "talent_collection"
}

func (h *TalentCollectionHandler) ReadAll(c echo.Context) error {
	records, err := h.TalentCollectionRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *TalentCollectionHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.TalentCollectionRepository.GetByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *TalentCollectionHandler) CreateOne(c echo.Context) error {
	params := new(repo.TalentCollectionParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.TalentCollectionRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *TalentCollectionHandler) UpdateByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.TalentCollectionParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	record, err := h.TalentCollectionRepository.Update(id, *params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *TalentCollectionHandler) DeleteOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.TalentCollectionRepository.DeleteByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
