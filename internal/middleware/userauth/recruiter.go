package userauth

import (
	"errors"
	"strings"

	"github.com/10hourlabs/tenlog"
	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	"github.com/10hourlabs/tentn/internal/middleware/globalctx"
	"github.com/labstack/echo/v4"
)

type RecruiterAuth struct {
	RoleAuther
}

func NewRecruiterAuth() RoleAuther {
	return &RecruiterAuth{}
}

func (auth RecruiterAuth) Authorize(u *ent.User, ctx echo.Context) error {
	if u.Role != userrole.Recruiter {
		return echo.ErrUnauthorized
	}
	if !u.Approved {
		tenlog.Error("user not approved", "user", u.ID)
		return errors.New("account not approved")
	}
	if !auth.isPathAllowed(ctx.Request().URL.Path) {
		return echo.ErrUnauthorized
	}
	return globalctx.SetCurrentRecruiterContext(ctx, u)
}

func (RecruiterAuth) isPathAllowed(path string) bool {
	return strings.Contains(path, "/v1/recruiter")
}
