package middleware

import (
	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type roleAuther interface {
	authorize(user *ent.User, c echo.Context) error
	isPathAllowed(path string) bool
}

var userRepo = repo.NewUserRepository()

var roleAuth = map[userrole.Role]roleAuther{
	userrole.Talent:    newTalentAuth(),
	userrole.Recruiter: newRecruiterAuth(),
	userrole.Developer: newDeveloperAuth(),
}

func AuthorizieUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			if !token.Valid {
				return echo.ErrUnauthorized
			}
			claims := token.Claims.(jwt.MapClaims)
			userID := claims["sub"].(string)
			if err := setCurrentUserContext(c, userID); err != nil {
				return echo.ErrUnauthorized
			}
			user := c.Get(oneword.CurrentUser).(*ent.User)
			if err := roleAuth[user.Role].authorize(user, c); err != nil {
				return err
			}
			return next(c)
		}
	}
}

func setCurrentUserContext(ctx echo.Context, userID string) error {
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
	ctx.Set(oneword.CurrentUser, user)
	return nil
}
