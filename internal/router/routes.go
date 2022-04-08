package router

import (
	"github.com/10hourlabs/tentn/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func DefineRoutes() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger()) // TODO: use a custom logger
	e.Use(middleware.Recover())
	e.GET("/", handler.IndexHandler)
	e.GET("/health", handler.HealthHandler)

	// TODO: Group routes by recruiters or talent
	e.GET("/talent/auth", handler.TalentLoginHanlder)
	e.GET("/oauth2/google/callback", handler.GoogleOauth2CallbackHandler)

	e.GET("/recruiter/auth", handler.RecruiterLoginHanlder)
	e.GET("/oauth2/slack/callback", handler.SlackOauth2CallbackHandler)

	// recruiter logout
	e.GET("/recruiter/auth/logout", handler.LogoutHandler)
	e = DefineV1Routes(e)
	return e
}
