package userauth

import (
	"strings"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	"github.com/10hourlabs/tentn/internal/middleware/globalctx"
	"github.com/labstack/echo/v4"
)

type TalentAuth struct {
	RoleAuther
}

func NewTalentAuth() RoleAuther {
	return &TalentAuth{}
}

func (auth TalentAuth) Authorize(u *ent.User, ctx echo.Context) error {
	if u.Role != userrole.Talent {
		return echo.ErrUnauthorized
	}
	if !u.Approved {
		return echo.ErrUnauthorized
	}
	if !auth.isPathAllowed(ctx.Request().URL.Path) {
		return echo.ErrUnauthorized
	}
	return globalctx.SetCurrentTalentContext(ctx, u.ID)
}

func (TalentAuth) isPathAllowed(path string) bool {
	return strings.Contains(path, "/v1/talent")
}
