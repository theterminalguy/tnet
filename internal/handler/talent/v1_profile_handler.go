package talent

import (
	"fmt"
	"net/http"

	"github.com/10hourlabs/tentn/internal/handler"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type V1TalentProfileHandler struct {
	TalentService    *service.TalentService
	TalentRepository *repo.TalentRepository
}

func NewV1TalentProfileHandler() *V1TalentProfileHandler {
	return &V1TalentProfileHandler{
		TalentService:    service.NewTalentService(),
		TalentRepository: repo.NewTalentRepository(),
	}
}

func (h *V1TalentProfileHandler) Search(c echo.Context) error {
	return nil
}

func (h *V1TalentProfileHandler) ReadAll(c echo.Context) error {
	return nil
}

func (h *V1TalentProfileHandler) ReadByID(c echo.Context) error {
	user, err := handler.GetCurrentUser(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	talent, err := h.TalentRepository.GetTalentByUserID(user.ID)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, talent)
}

func (h *V1TalentProfileHandler) CreateOne(c echo.Context) error {
	params := new(repo.TalentParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	user, err := handler.GetCurrentUser(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	a, err := h.TalentService.CreateProfile(user, params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, a)
}

func (h *V1TalentProfileHandler) UpdateByID(c echo.Context) error {
	params := new(repo.TalentParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	user, err := handler.GetCurrentUser(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	a, vldErrs := h.TalentService.UpdateProfile(user, params)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, a)
}

func (h *V1TalentProfileHandler) DeleteOne(c echo.Context) error {
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
