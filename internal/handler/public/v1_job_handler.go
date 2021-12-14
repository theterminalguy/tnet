package public

import (
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type V1PublicJobHandler struct {
	JobRepository *repo.JobRepository
}

func NewV1PublicJobHandler() *V1PublicJobHandler {
	return &V1PublicJobHandler{
		JobRepository: repo.NewJobRepository(),
	}
}

func (*V1PublicJobHandler) Search(c echo.Context) error {
	return nil
}

// ReadAlll return all jobs scopd to a paritcular talent
// i.e. all jobs for which the talent has a job application
func (h *V1PublicJobHandler) ReadAll(c echo.Context) error {
	// TODO: implement pagination
	// most likely cursor based
	jobs, err := h.JobRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, jobs)
}

// ReadByID same as above but filters the above by id.
// Should not return any jobs that are not scoped to the talent,
// even if the correct ID is provided
func (h *V1PublicJobHandler) ReadByID(c echo.Context) error {
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

func (*V1PublicJobHandler) CreateOne(c echo.Context) error {
	return nil
}

func (*V1PublicJobHandler) UpdateByID(c echo.Context) error {
	return nil
}

func (*V1PublicJobHandler) DeleteOne(c echo.Context) error {
	return nil
}
