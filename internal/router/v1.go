package router

import (
	"github.com/10hourlabs/tentn/internal/handler"
	"github.com/10hourlabs/tentn/internal/middleware"
	"github.com/labstack/echo/v4"
)

type v1Router struct {
	group       *echo.Group
	middlewares []echo.MiddlewareFunc
	handlers    []RouteHandler
}

// BuildRoutes creates all the routes for the given namespace
// all routes in this namespace require authentication
// you'd have to provide a JWT token to access the routes
func (v1 *v1Router) BuildRoutes() {
	v1.group.Use(v1.middlewares...)
	for _, h := range v1.handlers {
		h.Restify(v1.group)
	}
}

func DefineV1Routes(e *echo.Echo) {
	publicV1Router := &v1Router{
		group: e.Group("/v1/public"),
		handlers: []RouteHandler{
			{
				Path:        "jobs",
				Handler:     handler.NewJobHandler(),
				Middlewares: nil,
			},
		},
	}
	publicV1Router.BuildRoutes()

	talentRouter := &v1Router{
		group: e.Group("/v1/talent"),
		middlewares: []echo.MiddlewareFunc{
			middleware.JWTAuthenticate(),
		},
		handlers: []RouteHandler{
			{
				Path:        "jobs",
				Handler:     handler.NewJobHandler(),
				Middlewares: nil,
			},
		},
	}
	talentRouter.BuildRoutes()

	recruiterRouter := &v1Router{
		group: e.Group("/v1/recruiter"),
		middlewares: []echo.MiddlewareFunc{
			middleware.JWTAuthenticate(),
		},
		handlers: []RouteHandler{
			{
				Path:        "jobs",
				Handler:     handler.NewJobHandler(),
				Middlewares: nil,
			},
		},
	}
	recruiterRouter.BuildRoutes()
}
