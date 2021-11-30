package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/10hourlabs/tentn/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func DefineRoutes() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", func(c echo.Context) error {
		// TODO replace with documentation homepage
		return c.String(http.StatusOK, "Talent Network API version 0.0.1")
	})
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	})
	e.GET("/auth", handler.GoogleLoginHandler)
	e.GET("/oauth2/google/callback", handler.GoogleOauth2CallbackHandler)
	router := NewV1Router(e)
	router.CreateRoutes()
	return router.Engine()
}

func createRoutes(g *echo.Group, rh RouteHandler) {
	resourcePath := fmt.Sprintf("/%s", rh.Handler.ResourceName())
	byIDPath := fmt.Sprintf("%s/:uuid", resourcePath)

	// GET /resources
	g.GET(resourcePath, rh.Handler.ReadAll, rh.Middleware...)

	// GET /resources/:id
	g.GET(byIDPath, rh.Handler.ReadByID, rh.Middleware...)

	// POST /resources
	g.POST(resourcePath, rh.Handler.CreateOne, rh.Middleware...)

	// PUT /resources/:id
	g.PUT(byIDPath, rh.Handler.UpdateByID, rh.Middleware...)

	// DELETE /resources/:id
	g.DELETE(byIDPath, rh.Handler.DeleteOne, rh.Middleware...)
}

type RequestHandler interface {
	ResourceName() string

	ReadAll(c echo.Context) error
	ReadByID(c echo.Context) error

	CreateOne(c echo.Context) error

	UpdateByID(c echo.Context) error

	DeleteOne(c echo.Context) error
}

type RouteHandler struct {
	Handler    RequestHandler
	Middleware []echo.MiddlewareFunc
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

// createRoutes creates all the routes for the given namespace
// all routes in this namespace require authentication
// you'd have to provide a JWT token to access the routes
func (v1 *v1Router) CreateRoutes() {
	config := middleware.JWTConfig{
		SigningKey: []byte(os.Getenv("JWT_SIGNED_SECRET")),
	}
	g := v1.Namespace()
	g.Use(middleware.JWTWithConfig(config))
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
