package handler

import (
	"net/http"

	repo "github.com/theterminalguy/tentn/internal/repository"
	"github.com/theterminalguy/tentn/internal/service"
	"github.com/labstack/echo/v4"
)

type SlackInstallationHandler struct {
	SlackInstallRepo         *repo.SlackAppInstallRepository
	SlackAppUnInstallService *service.SlackAppUnInstallService
}

func NewSlackInstallationHandler() *SlackInstallationHandler {
	return &SlackInstallationHandler{
		SlackInstallRepo:         repo.NewSlackAppInstallRepository(),
		SlackAppUnInstallService: service.NewSlackAppUnInstallService(),
	}
}

func (*SlackInstallationHandler) ResourceName() string {
	return "slack_installation"
}

func (*SlackInstallationHandler) Search(c echo.Context) error {
	return nil
}

// ReadAll, this endpoint is an internal endpoint
// if we do choose to make it public, then we should consider scoping the request to the correct
// recruiter
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

func (h *SlackInstallationHandler) UpdateByID(c echo.Context) error {
	params := new(service.SlackAppUninstallParams)
	err := c.Bind(params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if params.Event.Type != "app_uninstalled" {
		return c.String(http.StatusBadRequest, "invalid event type")
	}
	err = h.SlackAppUnInstallService.UnInstall(params)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "ok")
}

func (SlackInstallationHandler) DeleteOne(c echo.Context) error {
	return nil
}
