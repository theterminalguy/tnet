package router

import (
	"github.com/10hourlabs/tentn/internal/handler"
	"github.com/10hourlabs/tentn/internal/handler/authserver"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func DefineRoutes() *echo.Echo {
	e := echo.New()
	// TODO: use a custom logger
	// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Format: "method=${method}, uri=${uri}, status=${status} ${latency_human}\n",
	// }))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", handler.IndexHandler)
	e.GET("/health", handler.HealthHandler)
	e.GET("/how-to-search", handler.HowToSearchHandler)

	// TODO: Group routes by recruiters or talent
	e.GET("/talent/auth", handler.TalentLoginHanlder)
	e.GET("/oauth2/google/callback", handler.GoogleOauth2CallbackHandler)

	e.GET("/recruiter/auth", handler.RecruiterLoginHanlder)
	e.GET("/oauth2/slack/callback", handler.SlackOauth2CallbackHandler)

	e.POST("/oauth2/client/register", authserver.Oauth2RegistrationHandler)
	e.POST("/oauth2/client/token", authserver.Oauth2AccessTokenHandler)
	e.POST("/oauth2/client/refresh", authserver.Oauth2RefreshTokenHandler)

	// recruiter logout
	e.GET("/recruiter/auth/logout", handler.LogoutHandler)
	e = DefineV1Routes(e)
	return e
}
