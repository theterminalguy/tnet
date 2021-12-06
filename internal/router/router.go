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

type HTTPMethod string

const (
	GET    HTTPMethod = http.MethodGet
	POST   HTTPMethod = http.MethodPost
	PUT    HTTPMethod = http.MethodPut
	DELETE HTTPMethod = http.MethodDelete
)

type RouteHandler struct {
	Path        string
	Only        []HTTPMethod
	Except      []HTTPMethod
	Handler     RequestHandler
	Middlewares []echo.MiddlewareFunc
}

func (rh *RouteHandler) Restify(g *echo.Group) {
	resourcePath := "/" + rh.Path
	resourceByIDPath := resourcePath + "/:uuid"
	endpoints := make(map[HTTPMethod]func())
	endpoints[GET] = func() {
		g.GET(resourcePath, rh.Handler.ReadAll, rh.Middlewares...)
		g.GET(resourceByIDPath, rh.Handler.ReadByID, rh.Middlewares...)
	}
	endpoints[POST] = func() {
		g.POST(resourcePath, rh.Handler.CreateOne, rh.Middlewares...)
	}
	endpoints[PUT] = func() {
		g.PUT(resourceByIDPath, rh.Handler.UpdateByID, rh.Middlewares...)
	}
	endpoints[DELETE] = func() {
		g.DELETE(resourceByIDPath, rh.Handler.DeleteOne, rh.Middlewares...)
	}
	if len(rh.Only) > 0 {
		for _, method := range rh.Only {
			endpoints[method]()
		}
		return
	}
	allMethods := []HTTPMethod{GET, POST, PUT, DELETE}
	diff := func(a, b []HTTPMethod) []HTTPMethod {
		mb := make(map[HTTPMethod]bool, len(b))
		for _, m := range b {
			mb[m] = true
		}
		var ms []HTTPMethod
		for _, m := range a {
			if _, ok := mb[m]; !ok {
				ms = append(ms, m)
			}
		}
		return ms
	}
	allowedMethods := diff(allMethods, rh.Except)
	for _, method := range allowedMethods {
		endpoints[method]()
	}
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
	DefineV1Routes(e)
	return e
}
