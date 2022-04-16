package handler

import (
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Oauth2TokenHandler struct {
	Oauth2TokenService    *service.Oauth2TokenService
	Oauth2TokenRepository *repo.Oauth2TokenRepository
}

func NewOauth2TokenHandler() *Oauth2TokenHandler {
	return &Oauth2TokenHandler{
		Oauth2TokenService:    service.NewOauth2TokenService(),
		Oauth2TokenRepository: repo.NewOauth2TokenRepository(),
	}
}

func (*Oauth2TokenHandler) ResourceName() string {
	return "oauth2_token"
}

func (h *Oauth2TokenHandler) ReadAll(c echo.Context) error {
	records, err := h.Oauth2TokenRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *Oauth2TokenHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.Oauth2TokenRepository.GetByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *Oauth2TokenHandler) CreateOne(c echo.Context) error {
	params := new(repo.Oauth2TokenParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.Oauth2TokenRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *Oauth2TokenHandler) UpdateByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.Oauth2TokenParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	record, err := h.Oauth2TokenRepository.Update(id, *params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *Oauth2TokenHandler) DeleteOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.Oauth2TokenRepository.DeleteByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
