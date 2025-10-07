package public

import (
	"net/http"

	"github.com/labstack/echo/v4"
	repo "github.com/theterminalguy/tnet/internal/repository"
)

type V1TimezoneHandler struct {
}

func NewV1TimezoneHandler() *V1TimezoneHandler {
	return &V1TimezoneHandler{}
}

func (*V1TimezoneHandler) Search(c echo.Context) error {
	return nil
}

func (h *V1TimezoneHandler) ReadAll(c echo.Context) error {
	records := repo.ReturnTimezoneData()
	return c.JSON(http.StatusOK, records)
}

func (h *V1TimezoneHandler) ReadByID(c echo.Context) error {
	return nil
}

func (*V1TimezoneHandler) CreateOne(c echo.Context) error {
	return nil
}

func (*V1TimezoneHandler) UpdateByID(c echo.Context) error {
	return nil
}

func (*V1TimezoneHandler) DeleteOne(c echo.Context) error {
	return nil
}
