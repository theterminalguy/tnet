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
}

// BuildRoutes creates all the routes for the given namespace
// all routes in this namespace require authentication
// you'd have to provide a JWT token to access the routes
func (v1 *v1Router) BuildRoutes() {
	v1.group.Use(middleware.JWTWithConfig(jwtConfig))
	for _, h := range v1.handlers {
		h.Restify(v1.group)
	}
}

func NewV1Router(e *echo.Echo) *v1Router {
	return &v1Router{
		group: e.Group("/v1"),
		handlers: []RouteHandler{
			{
				Path:        "jobs",
				Handler:     handler.NewJobHandler(),
				Middlewares: nil,
			},
			{
				Path:        "jobs/applications",
				Handler:     handler.NewJobApplicationHandler(),
				Middlewares: nil,
			},
			{
				Path:        "talents",
				Except:      []HTTPMethod{POST},
				Handler:     handler.NewTalentHandler(),
				Middlewares: nil,
			},
			{
				Path:        "talents/portfolio-links",
				Handler:     handler.NewPortfolioLinkHandler(),
				Middlewares: nil,
			},
			{
				Path:        "talents/skills",
				Handler:     handler.NewSkillHandler(),
				Middlewares: nil,
			},
			{
				Path:        "talents/work-experiences",
				Handler:     handler.NewWorkExperienceHandler(),
				Middlewares: nil,
			},
			{
				Path:        "talents/educations",
				Handler:     handler.NewEducationHandler(),
				Middlewares: nil,
			},
			{
				Path:        "talents/emergency-contacts",
				Handler:     handler.NewEmergencyContactHandler(),
				Middlewares: nil,
			},
			{
				Path:        "partners",
				Handler:     handler.NewPartnerHandler(),
				Middlewares: nil,
			},
			{
				Path:        "partners/missions",
				Handler:     handler.NewMissionHandler(),
				Middlewares: nil,
			},
		},
	}
}
