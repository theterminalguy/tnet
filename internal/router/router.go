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

type Request string

const (
	SEARCH       Request = "SEARCH"
	READ_ALL     Request = "READ_ALL"
	READ_BY_ID   Request = "READ_BY_ID"
	CREATE_ONE   Request = "CREATE_ONE"
	UPDATE_BY_ID Request = "UPDATE_BY_ID"
	DELETE_BY_ID Request = "DELETE_BY_ID"
)

var RequestsAllowedByDefault []Request = []Request{
	SEARCH,
	READ_ALL,
	READ_BY_ID,
	CREATE_ONE,
	UPDATE_BY_ID,
	DELETE_BY_ID,
}

type RouteHandler struct {
	Path        string
	Only        []Request
	Except      []Request
	Handler     RequestHandler
	Middlewares []echo.MiddlewareFunc
}

func (rh *RouteHandler) Restify(g *echo.Group) {
	resourcePath := "/" + rh.Path
	resourceByIDPath := resourcePath + "/:uuid"
	searchPath := resourceByIDPath + "/search"
	endpoints := make(map[Request]func())
	endpoints[READ_ALL] = func() {
		g.GET(resourcePath, rh.Handler.ReadAll, rh.Middlewares...)
	}
	endpoints[READ_BY_ID] = func() {
		g.GET(resourceByIDPath, rh.Handler.ReadByID, rh.Middlewares...)
	}
	endpoints[SEARCH] = func() {
		g.GET(searchPath, rh.Handler.Search, rh.Middlewares...)
	}
	endpoints[CREATE_ONE] = func() {
		g.POST(resourcePath, rh.Handler.CreateOne, rh.Middlewares...)
	}
	endpoints[UPDATE_BY_ID] = func() {
		g.PUT(resourceByIDPath, rh.Handler.UpdateByID, rh.Middlewares...)
	}
	endpoints[DELETE_BY_ID] = func() {
		g.DELETE(resourceByIDPath, rh.Handler.DeleteOne, rh.Middlewares...)
	}
	if len(rh.Only) > 0 {
		for _, method := range rh.Only {
			endpoints[method]()
		}
		return
	}
	diff := func(a, b []Request) []Request {
		mb := make(map[Request]bool, len(b))
		for _, m := range b {
			mb[m] = true
		}
		var ms []Request
		for _, m := range a {
			if _, ok := mb[m]; !ok {
				ms = append(ms, m)
			}
		}
		return ms
	}
	reqs := diff(RequestsAllowedByDefault, rh.Except)
	for _, req := range reqs {
		endpoints[req]()
	}
}
