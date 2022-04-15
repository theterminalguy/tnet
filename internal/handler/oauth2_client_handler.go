package handler

import (
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Oauth2ClientHandler struct {
	Oauth2ClientService    *service.Oauth2ClientService
	Oauth2ClientRepository *repo.Oauth2ClientRepository
}

func NewOauth2ClientHandler() *Oauth2ClientHandler {
	return &Oauth2ClientHandler{
		Oauth2ClientService:    service.NewOauth2ClientService(),
		Oauth2ClientRepository: repo.NewOauth2ClientRepository(),
	}
}

func (*Oauth2ClientHandler) ResourceName() string {
	return "oauth2_client"
}

func (h *Oauth2ClientHandler) ReadAll(c echo.Context) error {
	records, err := h.Oauth2ClientRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *Oauth2ClientHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.Oauth2ClientRepository.GetByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *Oauth2ClientHandler) CreateOne(c echo.Context) error {
	params := new(repo.Oauth2ClientParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.Oauth2ClientRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *Oauth2ClientHandler) UpdateByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.Oauth2ClientParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	record, err := h.Oauth2ClientRepository.Update(id, *params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *Oauth2ClientHandler) DeleteOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.Oauth2ClientRepository.DeleteByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func ClientRegisterationHandler(c echo.Context) error {
	params := new(repo.Oauth2ClientParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.Oauth2ClientRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}
