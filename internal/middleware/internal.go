package middleware

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/theterminalguy/tnet/internal/middleware/header"
	"github.com/theterminalguy/tnet/util/collection"
)

var allowedUsers = []string{
	"sp@theterminalguy.com",
	"dotun.oyelakin@theterminalguy.com",
	"abiodun.solomon@theterminalguy.com",
	"fortune.nwankwo@theterminalguy.com",
	"drey.olawaye@theterminalguy.com",
	"ahbot-slack@theterminalguy.com",
}

func AuthInternalRequest() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Request().Header.Get(header.X_TN_INTERNAL_USER_ID)
			password := c.Request().Header.Get(header.X_TN_INTERNAL_API_KEY)
			if userID == "" || password == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "user id or password missing")
			}
			if !collection.Contains(allowedUsers, userID) {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid user or password")
			}
			if password != os.Getenv("INTERNAL_API_KEY") {
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}
