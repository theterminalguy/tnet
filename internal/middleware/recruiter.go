package middleware

import (
	"strings"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	"github.com/labstack/echo/v4"
)

type recruiterAuth struct {
	roleAuther
}

func newRecruiterAuth() roleAuther {
	return &recruiterAuth{}
}

func (auth recruiterAuth) authorize(user *ent.User, ctx echo.Context) error {
	if user.Role != userrole.Recruiter {
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

func (recruiterAuth) isPathAllowed(path string) bool {
	return strings.Contains(path, "/v1/recruiter")
}
