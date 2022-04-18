package middleware

import (
	"strings"

	"github.com/10hourlabs/tentn/ent/schema/userrole"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func AuthenticateUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			if !token.Valid {
				return echo.ErrUnauthorized
			}
			claims := token.Claims.(jwt.MapClaims)
			currentUserRole := claims["role"]
			accountStatus := claims["approved"].(bool)
			currentPath := c.Request().URL.Path
			switch currentUserRole {
			case string(userrole.Recruiter):
				authenticateRecruiter(currentPath, accountStatus)
			case string(userrole.Talent):
				authenticateTalent(currentPath, accountStatus)
			case string(userrole.Developer):
				return next(c)
			default:
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}

func authenticateRecruiter(currentPath string, isApproved bool) error {
	// check if the path is in the list of paths that are allowed for recruiters
	if !isPathAllowedForRecruiter(currentPath) {
		return echo.ErrUnauthorized
	}
	// check if the recruiter has been approved
	if !isApproved {
		return echo.ErrUnauthorized
	}
	// TODO: prevent deleted accounts
	// This isn't done now as we are yet to figure out the best way to approach this
	// The obvious straightforward way is to check if the user is deleted, but that
	// would require a database query. We may have to cache the user data to allow for
	// faster access.
	return nil
}

func authenticateTalent(currentPath string, isApproved bool) error {
	// check if the path is in the list of paths that are allowed for talents
	if !isPathAllowedForTalent(currentPath) {
		return echo.ErrUnauthorized
	}
	// check if the talent has been approved
	if !isApproved {
		return echo.ErrUnauthorized
	}
	// TODO: prevent deleted accounts
	// This isn't done now as we are yet to figure out the best way to approach this
	// The obvious straightforward way is to check if the user is deleted, but that
	// would require a database query. We may have to cache the user data to allow for
	// faster access.
	return nil
}

func isPathAllowedForRecruiter(path string) bool {
	return strings.Index(path, "v1/recruiter") == 0
}

func isPathAllowedForTalent(path string) bool {
	return strings.Index(path, "v1/talent") == 0
}
