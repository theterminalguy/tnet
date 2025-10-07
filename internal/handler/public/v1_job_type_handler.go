package public

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/theterminalguy/tnet/ent/schema"
)

type V1JobTypeHandler struct {
}

func NewV1JobTypeHandler() *V1JobTypeHandler {
	return &V1JobTypeHandler{}
}

func (*V1JobTypeHandler) Search(c echo.Context) error {
	return nil
}

func (h *V1JobTypeHandler) ReadAll(c echo.Context) error {
	empTypes := schema.JobCategories()
	return c.JSON(http.StatusOK, empTypes)
}

func (h *V1JobTypeHandler) ReadByID(c echo.Context) error {
	return nil
}

func (*V1JobTypeHandler) CreateOne(c echo.Context) error {
	return nil
}

func (*V1JobTypeHandler) UpdateByID(c echo.Context) error {
	return nil
}

func (*V1JobTypeHandler) DeleteOne(c echo.Context) error {
	return nil
}
