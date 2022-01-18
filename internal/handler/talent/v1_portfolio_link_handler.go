package talent

import (
	"fmt"
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type V1PortfolioLinkHandler struct {
	PortfolioLinkService    *service.PortfolioLinkService
	PortfolioLinkRepository repo.PortfolioLinkQuerier
}

func NewV1PortfolioLinkHandler(pfLinkQuerier repo.PortfolioLinkQuerier) *V1PortfolioLinkHandler {
	return &V1PortfolioLinkHandler{
		PortfolioLinkService:    service.NewPortfolioLinkService(),
		PortfolioLinkRepository: pfLinkQuerier,
	}
}

func (h *V1PortfolioLinkHandler) Search(c echo.Context) error {
	return nil
}

func (h *V1PortfolioLinkHandler) ReadAll(c echo.Context) error {
	talent, err := GetCurrentTalent(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	records, err := h.PortfolioLinkRepository.GetAllForTalent(talent.ID)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *V1PortfolioLinkHandler) ReadByID(c echo.Context) error {
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

func (h *V1PortfolioLinkHandler) CreateOne(c echo.Context) error {
	err := VerifyTalentUUID(c)
	if err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
	}
	params := new(repo.PortfolioLinkParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.PortfolioLinkService.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *V1PortfolioLinkHandler) UpdateByID(c echo.Context) error {
	err := VerifyTalentUUID(c)
	if err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
	}
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

func (h *V1PortfolioLinkHandler) DeleteOne(c echo.Context) error {
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
