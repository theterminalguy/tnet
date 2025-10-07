package globalctx

import (
	"github.com/theterminalguy/tentn/ent"
	"github.com/theterminalguy/tentn/internal/repository/scope"
	"github.com/labstack/echo/v4"
)

func SetCurrentRecruiterContext(ctx echo.Context, u *ent.User) error {
	// TODO:
	// First try to get the u from cache
	// If the u is not in cache,
	// then get the u from database
	// and store the u in cache
	//
	// Suggested interim cache library:
	// https://github.com/allegro/bigcache
	currentRecruiter := scope.NewRecruiterScope(u)

	// We currently don't support multiple recruiter
	// I'm not sure if we need to support multiple recruiter
	// or not. This is open to discussion.
	// A possible solution would be to have one "recruiter" but
	// multiple "platform_users". For example,
	// one recruiter then multiple slack_app_users.
	ctx.Set("currentRecruiter", currentRecruiter)
	return nil
}
