package router

import (
	"os"

	"github.com/10hourlabs/tentn/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var jwtConfig middleware.JWTConfig = middleware.JWTConfig{
	SigningKey: []byte(os.Getenv("JWT_SIGNED_SECRET")),
}

type v1Router struct {
	group    *echo.Group
	handlers []RouteHandler
	//talentHandlers
	// recruiterHandlers
	// publicHandlers
	// ensure duplicate users don't signup
}

// BuildRoutes creates all the routes for the given namespace
// all routes in this namespace require authentication
// you'd have to provide a JWT token to access the routes
func (v1 *v1Router) BuildRoutes(m ...echo.MiddlewareFunc) {
	v1.group.Use(m...)
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
	publicV1Router.BuildRoutes(nil)

	protected := []echo.MiddlewareFunc{
		middleware.JWTWithConfig(jwtConfig),
	}
	talentRouter := &v1Router{
		group: e.Group("/v1/talent"),
		handlers: []RouteHandler{
			{
				Path:        "jobs",
				Handler:     handler.NewJobHandler(),
				Middlewares: nil,
			},
		},
	}
	talentRouter.BuildRoutes(protected...)
	
	recruiterRouter := &v1Router{
		group: e.Group("/v1/recruiter"),
		handlers: []RouteHandler{
			{
				Path:        "jobs",
				Handler:     handler.NewJobHandler(),
				Middlewares: nil,
			},
		},
	}
	recruiterRouter.BuildRoutes(protected...)
}
