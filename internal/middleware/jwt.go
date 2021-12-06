package middleware

import (
	"os"

	"github.com/labstack/echo/v4"
	mdlw "github.com/labstack/echo/v4/middleware"
)

var jwtConfig mdlw.JWTConfig = mdlw.JWTConfig{
	SigningKey: []byte(os.Getenv("JWT_SIGNED_SECRET")),
}

func JWTAuthenticate() echo.MiddlewareFunc {
	return mdlw.JWTWithConfig(jwtConfig)
}
