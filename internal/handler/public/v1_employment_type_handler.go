package public

import (
	"net/http"

	"github.com/theterminalguy/tentn/ent/schema"
	"github.com/labstack/echo/v4"
)

type V1EmploymentTypeHandler struct {
}

func NewV1EmploymentTypeHandler() *V1EmploymentTypeHandler {
	return &V1EmploymentTypeHandler{}
}

func (*V1EmploymentTypeHandler) Search(c echo.Context) error {
	return nil
}

func (h *V1EmploymentTypeHandler) ReadAll(c echo.Context) error {
	empTypes := schema.EmploymentTypes()
	return c.JSON(http.StatusOK, empTypes)
}

func (h *V1EmploymentTypeHandler) ReadByID(c echo.Context) error {
	return nil
}

func (*V1EmploymentTypeHandler) CreateOne(c echo.Context) error {
	return nil
}

func (*V1EmploymentTypeHandler) UpdateByID(c echo.Context) error {
	return nil
}

func (*V1EmploymentTypeHandler) DeleteOne(c echo.Context) error {
	return nil
}
