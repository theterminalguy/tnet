package handler

import (
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/internal/task"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type InternalTaskHandler struct {
	InternalTaskService    *service.InternalTaskService
	InternalTaskRepository *repo.InternalTaskRepository
}

func NewInternalTaskHandler() *InternalTaskHandler {
	return &InternalTaskHandler{
		InternalTaskService:    service.NewInternalTaskService(),
		InternalTaskRepository: repo.NewInternalTaskRepository(),
	}
}

func (*InternalTaskHandler) ResourceName() string {
	return "internal_task"
}

func (*InternalTaskHandler) Search(c echo.Context) error {
	return nil
}

func (h *InternalTaskHandler) ReadAll(c echo.Context) error {
	records, err := h.InternalTaskRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *InternalTaskHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.InternalTaskRepository.GetByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *InternalTaskHandler) CreateOne(c echo.Context) error {
	params := new(repo.InternalTaskParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.InternalTaskService.Create(task.Run, *params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"data": record,
			"err":  err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *InternalTaskHandler) UpdateByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.InternalTaskParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	record, err := h.InternalTaskRepository.Update(id, *params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *InternalTaskHandler) DeleteOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.InternalTaskRepository.DeleteByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
