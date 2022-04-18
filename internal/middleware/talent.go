package middleware

import (
	"strings"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	"github.com/labstack/echo/v4"
)

type talentAuth struct {
	roleAuther
}

func newTalentAuth() roleAuther {
	return &talentAuth{}
}

func (auth talentAuth) authorize(user *ent.User, ctx echo.Context) error {
	if user.Role != userrole.Talent {
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

func (talentAuth) isPathAllowed(path string) bool {
	return strings.Contains(path, "/v1/talent")
}
