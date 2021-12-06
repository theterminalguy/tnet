package router

import (
	"github.com/labstack/echo/v4"
)

type Router struct {
	group       *echo.Group
	middlewares []echo.MiddlewareFunc
	handlers    []RouteHandler
}

// BuildRoutes creates all the routes for the given namespace
// all routes in this namespace require authentication
// you'd have to provide a JWT token to access the routes
func (r *Router) BuildRoutes() {
	r.group.Use(r.middlewares...)
	for _, h := range r.handlers {
		h.Restify(r.group)
	}
}

type RequestHandler interface {
	Search(c echo.Context) error
	ReadAll(c echo.Context) error
	ReadByID(c echo.Context) error

	CreateOne(c echo.Context) error

	UpdateByID(c echo.Context) error

	DeleteOne(c echo.Context) error
}

type HTTPMethod string

const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
	GET_ID HTTPMethod = "GET_ID"
	SEARCH HTTPMethod = "SEARCH"
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
	searchPath := resourceByIDPath + "/search"
	endpoints := make(map[HTTPMethod]func())
	endpoints[GET] = func() {
		g.GET(resourcePath, rh.Handler.ReadAll, rh.Middlewares...)
	}
	endpoints[GET_ID] = func() {
		g.GET(resourceByIDPath, rh.Handler.ReadByID, rh.Middlewares...)
	}
	endpoints[SEARCH] = func() {
		g.GET(searchPath, rh.Handler.ReadByID, rh.Middlewares...)
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
	allMethods := []HTTPMethod{GET, POST, PUT, DELETE, GET_ID, SEARCH}
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
