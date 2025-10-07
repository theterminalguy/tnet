package talent

import (
	"fmt"
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/repository/scope"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type V1WorkExperienceHandler struct {
	WorkExperienceService    *service.WorkExperienceService
	WorkExperienceRepository repo.WorkExperienceQuerier
}

func NewV1WorkExperienceHandler(wrkExpQurier repo.WorkExperienceQuerier) *V1WorkExperienceHandler {
	return &V1WorkExperienceHandler{
		WorkExperienceService:    service.NewWorkExperienceService(),
		WorkExperienceRepository: wrkExpQurier,
	}
}

func (h *V1WorkExperienceHandler) Search(c echo.Context) error {
	return nil
}

func (h *V1WorkExperienceHandler) ReadAll(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	records, err := currentTalent.GetWorkExperiences()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *V1WorkExperienceHandler) ReadByID(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := currentTalent.GetWorkExperienceByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1WorkExperienceHandler) CreateOne(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	params := new(repo.WorkExperienceParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params.TalentID = currentTalent.GetID()
	record, err := h.WorkExperienceRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *V1WorkExperienceHandler) UpdateByID(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.WorkExperienceParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	record, vldErrs := currentTalent.UpdateWorkExperience(id, *params)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1WorkExperienceHandler) DeleteOne(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = currentTalent.DeleteWorkExperience(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
