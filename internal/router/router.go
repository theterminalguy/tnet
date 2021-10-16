package router

import (
	"fmt"

	"github.com/10hourlabs/tentn/internal/handler"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func DefineRoutes() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	return NewV1Router().createRoutes(e)
}

func createRoutes(namespace string, rh RouteHandler, e *echo.Echo) {
	basePath := fmt.Sprintf("/%s/%s", namespace, rh.Handler.ResourceName())
	allPath := fmt.Sprintf("%s", basePath)
	byIDPath := fmt.Sprintf("%s/:uuid", basePath)

	// GET /resources
	e.GET(allPath, rh.Handler.ReadAll, rh.Middleware...)

	// GET /resources/:id
	e.GET(byIDPath, rh.Handler.ReadByID, rh.Middleware...)

	// POST /resources
	e.POST(allPath, rh.Handler.CreateOne, rh.Middleware...)

	// PUT /resources/:id
	e.PUT(byIDPath, rh.Handler.UpdateByID, rh.Middleware...)

	// DELETE /resources/:id
	e.DELETE(byIDPath, rh.Handler.DeleteOne, rh.Middleware...)
}

type RequestHandler interface {
	ResourceName() string

	ReadAll(c echo.Context) error
	ReadByID(c echo.Context) error

	CreateOne(c echo.Context) error

	UpdateByID(c echo.Context) error

	DeleteOne(c echo.Context) error
}

type RouteHandler struct {
	Handler    RequestHandler
	Middleware []echo.MiddlewareFunc
}

type v1Router struct {
	namespace string
	handlers  []RouteHandler
}

func (v1 *v1Router) createRoutes(e *echo.Echo) *echo.Echo {
	for _, h := range v1.handlers {
		createRoutes(v1.namespace, h, e)
	}
	return e
}

func NewV1Router() *v1Router {
	// TODO: use Echo Group construct
	m := []echo.MiddlewareFunc{}
	return &v1Router{
		namespace: "v1",
		handlers: []RouteHandler{
			{handler.NewJobHandler(), m},
			{handler.NewApplicantHandler(), m},
			{handler.NewJobApplicationHandler(), m},
			{handler.NewPortfolioLinkHandler(), m},
			{handler.NewSkillHandler(), m},
			{handler.NewWorkExperienceHandler(), m},
		},
	}
}
