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
			isAccountApproved := claims["approved"].(bool)
			currentPath := c.Request().URL.Path
			switch currentUserRole {
			case string(userrole.Recruiter):
				authenticateRecruiter(isAccountApproved, currentPath)
			case string(userrole.Talent):
				authenticateTalent(isAccountApproved, currentPath)
			case string(userrole.Developer):
				return next(c)
			default:
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}

func authenticateRecruiter(isApproved bool, currentPath string) error {
	// check if the path is in the list of paths that are allowed for recruiters
	if !isPathAllowedForRecruiter(currentPath) {
		return echo.ErrUnauthorized
	}
	// check if the recruiter has been approved
	if !isApproved {
		return echo.ErrUnauthorized
	}
	return nil
}

func authenticateTalent(isApproved bool, currentPath string) error {
	// check if the path is in the list of paths that are allowed for talents
	if !isPathAllowedForTalent(currentPath) {
		return echo.ErrUnauthorized
	}
	// check if the talent has been approved
	if !isApproved {
		return echo.ErrUnauthorized
	}
	return nil
}

func isPathAllowedForRecruiter(path string) bool {
	return strings.Index(path, "v1/recruiter") == 0
}

func isPathAllowedForTalent(path string) bool {
	return strings.Index(path, "v1/talent") == 0
}
