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

	// define READ paths
	readAllPath := fmt.Sprintf("%s/ReadAll", basePath)
	readbyIDPath := fmt.Sprintf("%s/ReadByID/:uuid", basePath)

	// define CREATE paths
	createOnePath := fmt.Sprintf("%s/CreateOne", basePath)

	// define UPDATE paths
	updateByIDPath := fmt.Sprintf("%s/UpdateByID/:uuid", basePath)

	// define DELETE paths
	deleteOnePath := fmt.Sprintf("%s/DeleteOne/:uuid", basePath)

	e.GET(readAllPath, h.ReadAll, m...)
	e.GET(readbyIDPath, h.ReadByID, m...)

	e.POST(createOnePath, h.CreateOne, m...)

	e.PUT(updateByIDPath, h.UpdateByID, m...)

	e.DELETE(deleteOnePath, h.DeleteOne, m...)
}
