package userauth

import (
	"errors"
	"strings"

	"github.com/theterminalguy/tenlog"
	"github.com/theterminalguy/tentn/ent"
	"github.com/theterminalguy/tentn/ent/schema/userrole"
	"github.com/theterminalguy/tentn/internal/middleware/globalctx"
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
		tenlog.Error("user not approved", "user", u.ID)
		return errors.New("account not approved")
	}
	if !auth.isPathAllowed(ctx.Request().URL.Path) {
		return echo.ErrUnauthorized
	}
	return globalctx.SetCurrentTalentContext(ctx, u.ID)
}

func (TalentAuth) isPathAllowed(path string) bool {
	return strings.Contains(path, "/v1/talent")
}
