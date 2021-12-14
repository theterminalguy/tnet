package talent

import (
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type V1TaletJobHandler struct {
	JobRepository *repo.JobRepository
}

func NewV1TalentJobHandler() *V1TaletJobHandler {
	return &V1TaletJobHandler{
		JobRepository: repo.NewJobRepository(),
	}
}

func (*V1TaletJobHandler) Search(c echo.Context) error {
	return nil
}

func (h *V1TaletJobHandler) ReadAll(c echo.Context) error {
	// TODO: implement pagination
	// most likely cursor based
	jobs, err := h.JobRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, jobs)
}

func (h *V1TaletJobHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	j, err := h.JobRepository.GetByUUID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, j)
}

func (*V1TaletJobHandler) CreateOne(c echo.Context) error {
	return nil
}

func (*V1TaletJobHandler) UpdateByID(c echo.Context) error {
	return nil
}

func (*V1TaletJobHandler) DeleteOne(c echo.Context) error {
	return nil
}
