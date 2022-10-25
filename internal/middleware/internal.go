package middleware

import (
	"net/http"
	"os"

	"github.com/10hourlabs/tentn/internal/middleware/header"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/labstack/echo/v4"
)

var allowedUsers = []string{
	"sp@10hourlabs.com",
	"dotun.oyelakin@10hourlabs.com",
	"abiodun.solomon@10hourlabs.com",
	"fortune.nwankwo@10hourlabs.com",
	"drey.olawaye@10hourlabs.com",
	"ahbot-slack@10hourlabs.com",
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
