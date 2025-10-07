package talent

import (
	"fmt"
	"net/http"

	"github.com/theterminalguy/tentn/ent"
	repo "github.com/theterminalguy/tentn/internal/repository"
	"github.com/theterminalguy/tentn/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type V1TalentProfileHandler struct {
	TalentService    *service.TalentService
	TalentRepository repo.TalentQuerier
}

func NewV1TalentProfileHandler(talentQuerier repo.TalentQuerier) *V1TalentProfileHandler {
	return &V1TalentProfileHandler{
		TalentService:    service.NewTalentService(),
		TalentRepository: talentQuerier,
	}
}

func (h *V1TalentProfileHandler) Search(c echo.Context) error {
	return nil
}

func (h *V1TalentProfileHandler) ReadAll(c echo.Context) error {
	return nil
}

func (h *V1TalentProfileHandler) ReadByID(c echo.Context) error {
	user := c.Get("currentUser").(*ent.User)
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
	user := c.Get("currentUser").(*ent.User)
	a, err := h.TalentService.CreateProfile(user, *params)
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
	user := c.Get("currentUser").(*ent.User)
	a, vldErrs := h.TalentService.UpdateProfile(user, params)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, a)
}

func (h *V1TalentProfileHandler) DeleteOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.TalentRepository.DeleteByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
