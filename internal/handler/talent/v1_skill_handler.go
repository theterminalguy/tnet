package talent

import (
	"fmt"
	"net/http"

	repo "github.com/theterminalguy/tentn/internal/repository"
	"github.com/theterminalguy/tentn/internal/repository/scope"
	"github.com/theterminalguy/tentn/internal/search"
	"github.com/theterminalguy/tentn/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type V1SkillHandler struct {
	SkillService    *service.SkillService
	SkillRepository repo.SkillQuerier
}

func NewV1SkillHandler(skillQuerier repo.SkillQuerier) *V1SkillHandler {
	return &V1SkillHandler{
		SkillService:    service.NewSkillService(),
		SkillRepository: skillQuerier,
	}
}

func (h *V1SkillHandler) Search(c echo.Context) error {
	skillSearch := new(search.SkillSearch)
	query := c.QueryString()
	records, vldErrs := skillSearch.Search(query)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *V1SkillHandler) ReadAll(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	records, err := currentTalent.GetSkills()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *V1SkillHandler) ReadByID(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := currentTalent.GetSkillByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1SkillHandler) CreateOne(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	params := new(repo.SkillParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params.TalentID = currentTalent.GetID()
	record, err := h.SkillRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *V1SkillHandler) UpdateByID(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.SkillParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	record, vldErrs := currentTalent.UpdateSkill(id, *params)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1SkillHandler) DeleteOne(c echo.Context) error {
	currentTalent := c.Get("currentTalent").(*scope.TalentScope)
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = currentTalent.DeleteSkill(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
