package router

import (
	public_handler "github.com/10hourlabs/tentn/internal/handler/public"
	recruiter_handler "github.com/10hourlabs/tentn/internal/handler/recruiter"
	talent_handler "github.com/10hourlabs/tentn/internal/handler/talent"

	"github.com/10hourlabs/tentn/internal/middleware"
	"github.com/labstack/echo/v4"
)

func DefineV1Routes(e *echo.Echo) *echo.Echo {
	publicV1Router := &Router{
		group: e.Group("/v1/public"),
		handlers: []RouteHandler{
			{
				Path:        "jobs",
				Only:        []Request{READ_ALL},
				Handler:     public_handler.NewV1PublicJobHandler(),
				Middlewares: nil,
			},
			{
				Path:        "employment-types",
				Only:        []Request{READ_ALL},
				Handler:     public_handler.NewV1EmploymentTypeHandler(),
				Middlewares: nil,
			},
		},
	}
	publicV1Router.BuildRoutes()

	talentRouter := &Router{
		group: e.Group("/v1/talent"),
		middlewares: []echo.MiddlewareFunc{
			middleware.JWTAuthenticate(),
			middleware.EnforceTalent(),
		},
		handlers: []RouteHandler{
			{
				Path:        "jobs",
				Only:        []Request{READ_ALL},
				Handler:     talent_handler.NewV1TalentJobHandler(),
				Middlewares: nil,
			},
			{
				Path:        "skills",
				Except:      []Request{SEARCH},
				Handler:     talent_handler.NewV1SkillHandler(),
				Middlewares: nil,
			},
			{
				Path:        "profile",
				Except:      []Request{SEARCH},
				Handler:     talent_handler.NewV1TalentProfileHandler(),
				Middlewares: nil,
			},
		},
	}
	talentRouter.BuildRoutes()

	recruiterRouter := &Router{
		group: e.Group("/v1/recruiter"),
		middlewares: []echo.MiddlewareFunc{
			middleware.JWTAuthenticate(),
			middleware.EnforceApprovedRecruiter(),
		},
		handlers: []RouteHandler{
			{
				Path:        "jobs",
				Handler:     recruiter_handler.NewV1RecruiterJobHandler(),
				Middlewares: nil,
			},
		},
	}
	recruiterRouter.BuildRoutes()

	return e
}
