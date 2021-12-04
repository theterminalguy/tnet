package router

import (
	"fmt"
	"net/http"

	"github.com/10hourlabs/tentn/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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

func (rh *RouteHandler) Restify(g *echo.Group) {
	resourcePath := fmt.Sprintf("/%s", rh.Handler.ResourceName())
	byIDPath := fmt.Sprintf("%s/:uuid", resourcePath)

	// GET /resources
	g.GET(resourcePath, rh.Handler.ReadAll, rh.Middleware...)

	// GET /resources/:id
	g.GET(byIDPath, rh.Handler.ReadByID, rh.Middleware...)

	// POST /resources
	g.POST(resourcePath, rh.Handler.CreateOne, rh.Middleware...)

	// PUT /resources/:id
	g.PUT(byIDPath, rh.Handler.UpdateByID, rh.Middleware...)

	// DELETE /resources/:id
	g.DELETE(byIDPath, rh.Handler.DeleteOne, rh.Middleware...)
}

func DefineRoutes() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", func(c echo.Context) error {
		// TODO replace with documentation homepage
		return c.String(http.StatusOK, "Talent Network API version 0.0.1")
	})
	e.GET("/health", handler.HealthHandler)
	e.GET("/auth", handler.GoogleLoginHandler)
	e.GET("/oauth2/google/callback", handler.GoogleOauth2CallbackHandler)
	NewV1Router(e).BuildRoutes()
	return e
}
