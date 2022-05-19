package handler

import (
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SlackAppUserHandler struct {
	SlackAppUserService    *service.SlackAppUserService
	SlackAppUserRepository *repo.SlackAppUserRepository
}

func NewSlackAppUserHandler() *SlackAppUserHandler {
	return &SlackAppUserHandler{
		SlackAppUserService:    service.NewSlackAppUserService(),
		SlackAppUserRepository: repo.NewSlackAppUserRepository(),
	}
}

func (*SlackAppUserHandler) ResourceName() string {
	return "slack_app_user"
}

func (h *SlackAppUserHandler) ReadAll(c echo.Context) error {
	records, err := h.SlackAppUserRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *SlackAppUserHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.SlackAppUserRepository.GetByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *SlackAppUserHandler) CreateOne(c echo.Context) error {
	params := new(repo.SlackAppUserParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.SlackAppUserRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *SlackAppUserHandler) UpdateByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.SlackAppUserParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	record, err := h.SlackAppUserRepository.Update(id, *params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *SlackAppUserHandler) DeleteOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.SlackAppUserRepository.DeleteByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
