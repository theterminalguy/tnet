package middleware

import (
	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/google/uuid"
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
	// TODO:
	// First try to get the user from cache
	// If the user is not in cache,
	// then get the user from database
	// and store the user in cache
	//
	// Suggested interim cache library:
	// https://github.com/allegro/bigcache
	o2 := repo.NewOauth2ClientRepository()
	clientID := ctx.Get("clientID").(string)
	client, err := o2.GetByUUID(uuid.MustParse(clientID))
	if err != nil {
		return echo.ErrUnauthorized
	}
	if !client.Approved {
		return echo.ErrUnauthorized
	}
	if !auth.isPathAllowed(ctx.Request().URL.Path) {
		return echo.ErrUnauthorized
	}
	// TODO: cater for internal/external clients
	// There are certain things internal clients can do that external clients can't
	return nil
}

func (developerAuth) isPathAllowed(_ string) bool {
	// TODO: In the future,
	// there may be paths that are not allowed for developers

	// TODO: cater for internal/external clients
	// There are certain things internal clients can do that external clients can't
	return true
}
