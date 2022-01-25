package talent

import (
	"fmt"
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/search"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/oneword"
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
	talent, err := GetCurrentTalent(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	records, err := talent.GetSkills()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *V1SkillHandler) ReadByID(c echo.Context) error {
	talent, err := GetCurrentTalent(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := talent.GetSkillByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1SkillHandler) CreateOne(c echo.Context) error {
	talent, err := GetCurrentTalent(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.SkillParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params.TalentUUID = talent.Talent.UUID
	record, err := h.SkillRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *V1SkillHandler) UpdateByID(c echo.Context) error {
	talent, err := GetCurrentTalent(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.SkillParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	record, vldErrs := talent.UpdateSkill(id, *params)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1SkillHandler) DeleteOne(c echo.Context) error {
	talent, err := GetCurrentTalent(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = talent.DeleteSkill(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
