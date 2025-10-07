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

type V1EducationHandler struct {
	EducationService    *service.EducationService
	EducationRepository repo.EducationQuerier
}

func NewV1EducationHandler(eduQuerier repo.EducationQuerier) *V1EducationHandler {
	return &V1EducationHandler{
		EducationService:    service.NewEducationService(),
		EducationRepository: eduQuerier,
	}
}

func (h *V1EducationHandler) Search(c echo.Context) error {
	return nil
}

func (h *V1EducationHandler) ReadAll(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	records, err := currentTalent.GetEducations()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *V1EducationHandler) ReadByID(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := currentTalent.GetEducationByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1EducationHandler) CreateOne(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	params := new(repo.EducationParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params.TalentID = currentTalent.GetID()
	record, err := h.EducationRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *V1EducationHandler) UpdateByID(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.EducationParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	record, vldErrs := currentTalent.UpdateEducation(id, *params)
	if vldErrs != nil {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1EducationHandler) DeleteOne(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = currentTalent.DeleteEducation(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
