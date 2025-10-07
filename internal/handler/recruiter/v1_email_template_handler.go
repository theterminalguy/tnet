package recruiter

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	repo "github.com/theterminalguy/tnet/internal/repository"
	"github.com/theterminalguy/tnet/internal/repository/scope"
	"github.com/theterminalguy/tnet/internal/service"
)

type V1RecruiterEmailTemplateHandler struct {
	EmailTemplateService    *service.EmailTemplateService
	EmailTemplateRepository *repo.EmailTemplateRepository
}

func NewV1RecruiterEmailTemplateHandler() *V1RecruiterEmailTemplateHandler {
	return &V1RecruiterEmailTemplateHandler{
		EmailTemplateService:    service.NewEmailTemplateService(),
		EmailTemplateRepository: repo.NewEmailTemplateRepository(),
	}
}

func (h *V1RecruiterEmailTemplateHandler) Search(c echo.Context) error {
	return nil
}

func (h *V1RecruiterEmailTemplateHandler) ReadAll(c echo.Context) error {
	records, err := h.EmailTemplateRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *V1RecruiterEmailTemplateHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.EmailTemplateRepository.GetByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1RecruiterEmailTemplateHandler) CreateOne(c echo.Context) error {
	currentRecruiter := c.Get("currentRecruiter").(*scope.RecruiterScope)
	params := new(repo.EmailTemplateParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params.UserID = currentRecruiter.GetID()
	record, err := h.EmailTemplateRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *V1RecruiterEmailTemplateHandler) UpdateByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.EmailTemplateParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	record, err := h.EmailTemplateRepository.Update(id, *params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1RecruiterEmailTemplateHandler) DeleteOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.EmailTemplateRepository.DeleteByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
