package router

import (
	public_handler "github.com/10hourlabs/tentn/internal/handler/public"
	recruiter_handler "github.com/10hourlabs/tentn/internal/handler/recruiter"
	talent_handler "github.com/10hourlabs/tentn/internal/handler/talent"
	repo "github.com/10hourlabs/tentn/internal/repository"

	"github.com/10hourlabs/tentn/internal/middleware"
	"github.com/labstack/echo/v4"
)

func DefineV1Routes(e *echo.Echo) *echo.Echo {
	// TODOS: Look into what route verb should be allowed for each endpoint.
	publicV1Router := &Router{
		group: e.Group("/v1/public"),
		handlers: []RouteHandler{
			{
				Path:        "jobs",
				Only:        []Request{READ_ALL},
				Handler:     public_handler.NewV1PublicJobHandler(repo.NewJobRepository()),
				Middlewares: nil,
			},
			{
				Path:        "employment-types",
				Only:        []Request{READ_ALL},
				Handler:     public_handler.NewV1EmploymentTypeHandler(),
				Middlewares: nil,
			},
			{
				Path:        "job-types",
				Only:        []Request{READ_ALL},
				Handler:     public_handler.NewV1JobTypeHandler(),
				Middlewares: nil,
			},
			{
				Path:        "timezones",
				Only:        []Request{READ_ALL},
				Handler:     public_handler.NewV1TimezoneHandler(),
				Middlewares: nil,
			},
		},
	}
	publicV1Router.BuildRoutes()

	talentRouter := &Router{
		group: e.Group("/v1/talent"),
		middlewares: []echo.MiddlewareFunc{
			middleware.ExtractJWTTokenFromWebSession(),
			middleware.JWTAuthenticate(),
			middleware.AuthorizieUser(),
		},
		handlers: []RouteHandler{
			{
				Path:        "skills",
				Except:      []Request{SEARCH},
				Handler:     talent_handler.NewV1SkillHandler(repo.NewSkillRepository()),
				Middlewares: nil,
			},
			{
				SingularPath: "profile",
				Except:       []Request{SEARCH},
				Handler:      talent_handler.NewV1TalentProfileHandler(repo.NewTalentRepository()),
				Middlewares:  nil,
			},
			{
				Path:        "portfolio-links",
				Except:      []Request{SEARCH, READ, UPDATE},
				Handler:     talent_handler.NewV1PortfolioLinkHandler(repo.NewPortfolioLinkRepository()),
				Middlewares: nil,
			},
			{
				Path:        "work-experiences",
				Except:      []Request{SEARCH, READ, UPDATE},
				Handler:     talent_handler.NewV1WorkExperienceHandler(repo.NewWorkExperienceRepository()),
				Middlewares: nil,
			},
			{
				Path:        "educations",
				Except:      []Request{SEARCH, READ, UPDATE},
				Handler:     talent_handler.NewV1EducationHandler(repo.NewEducationRepository()),
				Middlewares: nil,
			},
			{
				Path:        "job-applications",
				Except:      []Request{READ, UPDATE, UPDATE_BY_ID, DELETE_BY_ID},
				Handler:     talent_handler.NewV1JobApplicationHandler(repo.NewJobApplicationRepository()),
				Middlewares: nil,
			},
		},
	}
	talentRouter.BuildRoutes()

	recruiterRouter := &Router{
		group: e.Group("/v1/recruiter"),
		middlewares: []echo.MiddlewareFunc{
			middleware.ExtractJWTTokenFromWebSession(),
			middleware.JWTAuthenticate(),
			middleware.AuthorizieUser(),
		},
		handlers: []RouteHandler{
			{
				Path:        "jobs",
				Handler:     recruiter_handler.NewV1RecruiterJobHandler(repo.NewJobRepository()),
				Middlewares: nil,
			},
			{
				Path:        "talents",
				Handler:     recruiter_handler.NewV1TalentSearchFilterHandler(),
				Middlewares: nil,
				Only:        []Request{SEARCH, READ_BY_ID},
			},
			{
				Path:        "job-applications",
				Handler:     recruiter_handler.NewV1RecruiterJobApplicationHandler(repo.NewJobApplicationRepository()),
				Middlewares: nil,
			},
			{
				Path:        "email-templates",
				Handler:     recruiter_handler.NewV1RecruiterEmailTemplateHandler(),
				Middlewares: nil,
			},
			{
				Path:        "talent-collections",
				Handler:     recruiter_handler.NewV1TalentCollectionHandler(),
				Middlewares: nil,
			},
		},
	}
	recruiterRouter.BuildRoutes()

	return e
}
