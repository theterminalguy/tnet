package middleware

import (
	"os"
	"strings"

	"github.com/10hourlabs/tenlog"
	"github.com/10hourlabs/tentn/internal/entsm"
	"github.com/10hourlabs/tentn/internal/tokgen"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ExtractJWTTokenFromWebSession extracts the JWT token from the session
func ExtractJWTTokenFromWebSession() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			store := entsm.GetSessionStore()
			session, err := store.Get(c.Request(), entsm.DefaultSessionName)
			if _, ok := session.Values["authenticated"].(bool); ok && err == nil {
				// if the request is not coming from the app's host,
				// then the session is not valid
				if strings.Index(c.Request().Host, os.Getenv("APP_HOST")) == -1 {
					// TODO: we should add more support to prevent
					// CSRF attacks. Since any client can fake the host,
					// we should add more checks here. Like a specifical csrf token
					// that is sent when the page is loaded.
					// https://portswigger.net/web-security/csrf/tokens
					// https://github.com/gorilla/csrf
					return echo.ErrUnauthorized
				}
				// TODO: how does this handle expired sessions?
				tok := session.Values["token"].(string)
				c.Request().Header.Set("Authorization", "Bearer "+tok)
			}
			if err != nil {
				tenlog.Error("Failed to get session: ", err)
			}
			return next(c)
		}
	}
}

// ValidateJWT verifies that the JWT token is valid, meaning:
// 1. The token is signed by the correct key
// 2. The token is not expired
// 3. The token contains the correct claims
func ValidateJWT() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: tokgen.DefaultSigningKey,
		Claims:     &tokgen.JWTClaims{},
	})
}
