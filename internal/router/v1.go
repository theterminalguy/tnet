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
	engine   *echo.Echo // ideally, the type for engine this should be an interface
	handlers []RouteHandler
}

func (v1 *v1Router) Engine() *echo.Echo {
	return v1.engine
}

func (v1 *v1Router) Namespace() *echo.Group {
	return v1.Engine().Group("/v1")
}

// CreateRoutes creates all the routes for the given namespace
// all routes in this namespace require authentication
// you'd have to provide a JWT token to access the routes
func (v1 *v1Router) CreateRoutes() {
	g := v1.Namespace()
	g.Use(middleware.JWTWithConfig(jwtConfig))
	for _, h := range v1.handlers {
		createRoutes(g, h)
	}
}

func NewV1Router(e *echo.Echo) *v1Router {
	return &v1Router{
		engine: e,
		handlers: []RouteHandler{
			{handler.NewJobHandler(), nil},
			{handler.NewJobApplicationHandler(), nil},
			{handler.NewPortfolioLinkHandler(), nil},
			{handler.NewSkillHandler(), nil},
			{handler.NewWorkExperienceHandler(), nil},
			{handler.NewEducationHandler(), nil},
			{handler.NewEmergencyContactHandler(), nil},
			{handler.NewTalentHandler(), nil},
			{handler.NewPartnerHandler(), nil},
			{handler.NewMissionHandler(), nil},
		},
	}
}
