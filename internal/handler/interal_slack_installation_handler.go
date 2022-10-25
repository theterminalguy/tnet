package handler

import (
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/labstack/echo/v4"
)

type SlackInstallationHandler struct {
	SlackInstallRepo *repo.SlackAppInstallRepository
}

func NewSlackInstallationHandler() *SlackInstallationHandler {
	return &SlackInstallationHandler{
		SlackInstallRepo: repo.NewSlackAppInstallRepository(),
	}
}

func (*SlackInstallationHandler) ResourceName() string {
	return "slack_installation"
}

func (*SlackInstallationHandler) Search(c echo.Context) error {
	return nil
}

func (h *SlackInstallationHandler) ReadAll(c echo.Context) error {

	teamID := c.QueryParam("team_id")
	if teamID != "" {
		record, err := h.SlackInstallRepo.GetByTeamID(teamID)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, record)
	}
	records, err := h.SlackInstallRepo.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (SlackInstallationHandler) ReadByID(c echo.Context) error {
	return nil
}

func (SlackInstallationHandler) CreateOne(c echo.Context) error {
	return nil
}

func (SlackInstallationHandler) UpdateByID(c echo.Context) error {
	return nil
}

func (SlackInstallationHandler) DeleteOne(c echo.Context) error {
	return nil
}
