package router

import (
	"fmt"

	"github.com/10hourlabs/tentn/internal/handler"
	"github.com/labstack/echo"
)

type RouteHandler interface {
	BasePath() string

	ReadByID(c echo.Context) error
	ReadAll(c echo.Context) error

	CreateOne(c echo.Context) error

	UpdateByID(c echo.Context) error

	DeleteOne(c echo.Context) error
}

func DefineRoutes() {
	e := echo.New()

	// JobController Routes
	createRoutes(handler.JobHandler{}, e)
}

func createRoutes(h RouteHandler, e *echo.Echo, m ...echo.MiddlewareFunc) {
	basePath := fmt.Sprintf("/%s", h.BasePath())

	e.GET(fmt.Sprintf("%s/ReadByID", basePath), h.ReadByID, m...)
	e.GET(fmt.Sprintf("%s/ReadAll", basePath), h.ReadAll, m...)

	e.POST(fmt.Sprintf("%s/CreateOne", basePath), h.CreateOne, m...)

	e.PUT(fmt.Sprintf("%s/UpdateByID", basePath), h.UpdateByID, m...)

	e.DELETE(fmt.Sprintf("%s/DeleteOne", basePath), h.DeleteOne, m...)
}
