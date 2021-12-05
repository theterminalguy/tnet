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

type PortfolioLinkHandler struct {
	PortfolioLinkService    *service.PortfolioLinkService
	PortfolioLinkRepository *repo.PortfolioLinkRepository
}

func NewPortfolioLinkHandler() *PortfolioLinkHandler {
	return &PortfolioLinkHandler{
		PortfolioLinkService:    service.NewPortfolioLinkService(),
		PortfolioLinkRepository: repo.NewPortfolioLinkRepository(),
	}
}

func (h *PortfolioLinkHandler) Search(c echo.Context) error {
	pfLinkSearch := new(search.PortfolioLinkSearch)
	query := c.QueryString()
	records, vldErrs := pfLinkSearch.Search(query)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *PortfolioLinkHandler) ReadAll(c echo.Context) error {
	records, err := h.PortfolioLinkRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *PortfolioLinkHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.PortfolioLinkRepository.GetByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *PortfolioLinkHandler) CreateOne(c echo.Context) error {
	params := new(repo.PortfolioLinkParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.PortfolioLinkRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *PortfolioLinkHandler) UpdateByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.PortfolioLinkParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	record, vldErrs := h.PortfolioLinkRepository.Update(id, *params)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *PortfolioLinkHandler) DeleteOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.PortfolioLinkRepository.DeleteByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
