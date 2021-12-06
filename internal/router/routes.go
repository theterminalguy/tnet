package router

import (
	"net/http"

	"github.com/10hourlabs/tentn/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func DefineRoutes() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger()) // TODO: use a custom logger
	e.Use(middleware.Recover())
	e.GET("/", func(c echo.Context) error {
		// TODO replace with documentation homepage
		return c.String(http.StatusOK, "Talent Network API version 0.0.1")
	})
	e.GET("/health", handler.HealthHandler)
	e.GET("/talent/auth", handler.TalentLoginHanlder)
	e.GET("/oauth2/google/callback", handler.GoogleOauth2CallbackHandler)
	e = DefineV1Routes(e)
	return e
}
