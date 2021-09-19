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

	// define READ paths
	readbyIDPath := fmt.Sprintf("%s/ReadByID", basePath)
	readAllPath := fmt.Sprintf("%s/ReadAll", basePath)

	// define CREATE paths
	createOnePath := fmt.Sprintf("%s/CreateOne", basePath)

	// define UPDATE paths
	updateByIDPath := fmt.Sprintf("%s/UpdateByID", basePath)

	// define DELETE paths
	deleteOnePath := fmt.Sprintf("%s/DeleteOne", basePath)

	e.GET(readbyIDPath, h.ReadByID, m...)
	e.GET(readAllPath, h.ReadAll, m...)

	e.POST(createOnePath, h.CreateOne, m...)

	e.PUT(updateByIDPath, h.UpdateByID, m...)

	e.DELETE(deleteOnePath, h.DeleteOne, m...)
}
