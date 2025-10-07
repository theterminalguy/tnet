package talent

import (
	"fmt"
	"net/http"

	repo "github.com/theterminalguy/tentn/internal/repository"
	"github.com/theterminalguy/tentn/internal/repository/scope"
	"github.com/theterminalguy/tentn/internal/service"
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
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	records, err := currentTalent.GetPortfolioLinks()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *V1PortfolioLinkHandler) ReadByID(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := currentTalent.GetPortfolioLinkByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1PortfolioLinkHandler) CreateOne(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	params := new(repo.PortfolioLinkParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params.TalentID = currentTalent.GetID()
	record, err := h.PortfolioLinkService.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *V1PortfolioLinkHandler) UpdateByID(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.PortfolioLinkParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	record, vldErrs := currentTalent.UpdatePortfolioLink(id, *params)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1PortfolioLinkHandler) DeleteOne(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = currentTalent.DeletePortfolioLink(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
