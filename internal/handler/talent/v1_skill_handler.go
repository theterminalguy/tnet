package talent

import (
	"fmt"
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type V1SkillHandler struct {
	SkillService    *service.SkillService
	SkillRepository *repo.SkillRepository
}

func NewV1SkillHandler() *V1SkillHandler {
	return &V1SkillHandler{
		SkillService:    service.NewSkillService(),
		SkillRepository: repo.NewSkillRepository(),
	}
}

func (*V1SkillHandler) Search(c echo.Context) error {
	return nil
}

func (h *V1SkillHandler) ReadAll(c echo.Context) error {
	records, err := h.SkillRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *V1SkillHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.SkillRepository.GetByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1SkillHandler) CreateOne(c echo.Context) error {
	params := new(repo.SkillParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.SkillRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *V1SkillHandler) UpdateByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.SkillParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	record, vldErrs := h.SkillRepository.Update(id, *params)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1SkillHandler) DeleteOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.SkillRepository.DeleteByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
