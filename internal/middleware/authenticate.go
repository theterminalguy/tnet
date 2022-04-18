package middleware

import (
	"strings"

	"github.com/10hourlabs/tenlog"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var userRepo = repo.NewUserRepository()

func AuthenticateUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			if !token.Valid {
				return echo.ErrUnauthorized
			}
			claims := token.Claims.(jwt.MapClaims)
			userID := claims["sub"].(string)
			userRole := claims["role"].(string)
			isApproved := claims["approved"].(bool)
			currentPath := c.Request().URL.Path
			switch userRole {
			case string(userrole.Recruiter):
				authenticateRecruiter(currentPath, isApproved)
			case string(userrole.Talent):
				authenticateTalent(currentPath, isApproved)
			case string(userrole.Developer):
				authenticateDeveloper(currentPath, isApproved)
			default:
				return echo.ErrUnauthorized
			}
			err := setCurrentUserContext(c, userID, userRole)
			if err != nil {
				tenlog.Error("Failed to set current user context", err)
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}

func setCurrentUserContext(ctx echo.Context, userID, userRole string) error {
	// TODO:
	// First try to get the user from cache
	// If the user is not in cache,
	// then get the user from database
	// and store the user in cache
	//
	// Suggested interim cache library:
	// https://github.com/allegro/bigcache
	user, err := userRepo.GetByID(uuid.MustParse(userID))
	if err != nil {
		return err
	}
	ctx.Set("currentUser", user)
	return nil
}

func authenticateRecruiter(path string, isApproved bool) error {
	if !isPathAllowedForRecruiter(path) {
		return echo.ErrUnauthorized
	}
	if !isApproved {
		return echo.ErrUnauthorized
	}
	return nil
}

func authenticateTalent(path string, isApproved bool) error {
	if !isPathAllowedForTalent(path) {
		return echo.ErrUnauthorized
	}
	if !isApproved {
		return echo.ErrUnauthorized
	}
	return nil
}

func authenticateDeveloper(path string, isApproved bool) error {
	if !isPathAllowedForDeveloper(path) {
		return echo.ErrUnauthorized
	}
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

func isPathAllowedForDeveloper(path string) bool {
	// TODO: In the future,
	// there may be paths that are not allowed for developers
	return true
}
