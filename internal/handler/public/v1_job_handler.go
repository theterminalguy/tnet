package public

import (
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type V1PublicJobHandler struct {
	JobRepository repo.JobQuerier
}

func NewV1PublicJobHandler(jobQuerier repo.JobQuerier) *V1PublicJobHandler {
	return &V1PublicJobHandler{
		JobRepository: jobQuerier,
	}
}

func (*V1PublicJobHandler) Search(c echo.Context) error {
	return nil
}

// ReadAlll return all jobs scopd to a paritcular talent
// i.e. all jobs for which the talent has a job application
func (h *V1PublicJobHandler) ReadAll(c echo.Context) error {
	jobs, err := h.JobRepository.GetAll(c.QueryParam("cursor"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, jobs)
}

// ReadByID same as above but filters the above by id.
// Should not return any jobs that are not scoped to the talent,
// even if the correct ID is provided
func (h *V1PublicJobHandler) ReadByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	j, err := h.JobRepository.GetByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, j)
}

func (*V1PublicJobHandler) CreateOne(c echo.Context) error {
	return nil
}

func (*V1PublicJobHandler) UpdateByID(c echo.Context) error {
	return nil
}

func (*V1PublicJobHandler) DeleteOne(c echo.Context) error {
	return nil
}
