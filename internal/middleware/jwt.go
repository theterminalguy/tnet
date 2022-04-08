package middleware

import (
	"os"

	"github.com/10hourlabs/tentn/internal/entsm"
	"github.com/labstack/echo/v4"
	mdlw "github.com/labstack/echo/v4/middleware"
)

var jwtConfig mdlw.JWTConfig = mdlw.JWTConfig{
	SigningKey: []byte(os.Getenv("JWT_SIGNED_SECRET")),
}

func JWTAuthenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			store := entsm.GetSessionStore()
			session, _ := store.Get(c.Request(), "tentn-session")
			if auth, ok := session.Values["authenticated"].(bool); ok || auth {
				var bearer = "Bearer " + session.Values["token"].(string)
				c.Request().Header.Set("Authorization", bearer)
			}
			return next(c)
		}
	}
}
