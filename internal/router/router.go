package router

import (
	"fmt"

	"github.com/10hourlabs/tentn/internal/handler"
	"github.com/labstack/echo"
)

type RouteHandler interface {
	ResourceName() string

	ReadAll(c echo.Context) error
	ReadByID(c echo.Context) error

	CreateOne(c echo.Context) error

	UpdateByID(c echo.Context) error

	DeleteOne(c echo.Context) error
}

func DefineRoutes() *echo.Echo {
	e := echo.New()

	// JobController Routes
	createRoutes(handler.NewJobHandler(), e)

	return e
}

func createRoutes(h RouteHandler, e *echo.Echo, m ...echo.MiddlewareFunc) {
	// TODO: use Echo Group construct
	namespace := "v1"

	basePath := fmt.Sprintf("/%s/%s", namespace, h.ResourceName())
	allPath := fmt.Sprintf("%s", basePath)
	byIDPath := fmt.Sprintf("%s/:uuid", basePath)

	// GET /resources
	e.GET(allPath, h.ReadAll, m...)

	// GET /resources/:id
	e.GET(byIDPath, h.ReadByID, m...)

	// POST /resources
	e.POST(byIDPath, h.CreateOne, m...)

	// PUT /resources/:id
	e.PUT(byIDPath, h.UpdateByID, m...)

	// DELETE /resources/:id
	e.DELETE(byIDPath, h.DeleteOne, m...)
}
