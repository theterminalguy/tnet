package router

import (
	"net/http"

	"github.com/10hourlabs/tentn/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type RequestHandler interface {
	ReadAll(c echo.Context) error
	ReadByID(c echo.Context) error

	CreateOne(c echo.Context) error

	UpdateByID(c echo.Context) error

	DeleteOne(c echo.Context) error
}

type RouteHandler struct {
	Path        string
	Handler     RequestHandler
	Middlewares []echo.MiddlewareFunc
}

func (rh *RouteHandler) Restify(g *echo.Group) {
	resourcePath := "/" + rh.Path
	resourceByIDPath := resourcePath + "/:uuid"
	g.GET(resourcePath, rh.Handler.ReadAll, rh.Middlewares...)
	g.POST(resourcePath, rh.Handler.CreateOne, rh.Middlewares...)
	g.GET(resourceByIDPath, rh.Handler.ReadByID, rh.Middlewares...)
	g.PUT(resourceByIDPath, rh.Handler.UpdateByID, rh.Middlewares...)
	g.DELETE(resourceByIDPath, rh.Handler.DeleteOne, rh.Middlewares...)
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
