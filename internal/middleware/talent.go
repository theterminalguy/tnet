package middleware

import (
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func EnforceTalent() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			if claims["role"] != string(userrole.Talent) {
				return c.JSON(401, map[string]interface{}{
					"message": "You are not a talent",
				})
			}
			return next(c)
		}
	}
}
