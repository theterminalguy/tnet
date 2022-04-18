package middleware

import (
	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	"github.com/labstack/echo/v4"
)

type developerAuth struct {
	roleAuther
}

func newDeveloperAuth() roleAuther {
	return &developerAuth{}
}

func (auth developerAuth) authorize(user *ent.User, ctx echo.Context) error {
	if user.Role != userrole.Developer {
		return echo.ErrUnauthorized
	}
	if !user.Approved {
		return echo.ErrUnauthorized
	}
	if !auth.isPathAllowed(ctx.Request().URL.Path) {
		return echo.ErrUnauthorized
	}
	return nil
}

func (developerAuth) isPathAllowed(_ string) bool {
	// TODO: In the future,
	// there may be paths that are not allowed for developers
	return true
}
