package router

import (
	"fmt"

	"github.com/10hourlabs/tentn/internal/handler"
	"github.com/labstack/echo"
)

type RouteHandler interface {
	BasePath() string

	ReadAll(c echo.Context) error
	ReadByID(c echo.Context) error

	CreateOne(c echo.Context) error

	UpdateByID(c echo.Context) error

	DeleteOne(c echo.Context) error
}

func DefineRoutes() *echo.Echo {
	e := echo.New()

	// JobController Routes
	createRoutes(handler.JobHandler{}, e)

	return e
}

func createRoutes(h RouteHandler, e *echo.Echo, m ...echo.MiddlewareFunc) {
	// TODO: use Echo Group construct
	namespace := "v1"

	basePath := fmt.Sprintf("/%s/%s", namespace, h.BasePath())

	// define READ paths
	readAllPath := fmt.Sprintf("%s/ReadAll", basePath)
	readbyIDPath := fmt.Sprintf("%s/ReadByID", basePath)

	// define CREATE paths
	createOnePath := fmt.Sprintf("%s/CreateOne", basePath)

	// define UPDATE paths
	updateByIDPath := fmt.Sprintf("%s/UpdateByID", basePath)

	// define DELETE paths
	deleteOnePath := fmt.Sprintf("%s/DeleteOne", basePath)

	e.GET(readAllPath, h.ReadAll, m...)
	e.GET(readbyIDPath, h.ReadByID, m...)

	e.POST(createOnePath, h.CreateOne, m...)

	e.PUT(updateByIDPath, h.UpdateByID, m...)

	e.DELETE(deleteOnePath, h.DeleteOne, m...)
}
