package middleware

import (
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func EnforceApprovedRecruiter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			if claims["role"] != string(userrole.Recruiter) {
				return echo.ErrUnauthorized
			}
			if claims["approved"] == false {
				return echo.ErrUnauthorized
			}
			// TODO: prevent deleted accounts
			// This isn't done now as we are yet to figure out the best way to approach this
			// The obvious straightforward way is to check if the user is deleted, but that
			// would require a database query. We may have to cache the user data to allow for
			// faster access.
			return next(c)
		}
	}
}
