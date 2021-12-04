package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

//HealtHanlder handles request that reaches the health endpoint
func HealthHandler(c echo.Context) error {
	return c.String(http.StatusOK, http.StatusText(http.StatusOK))
}
