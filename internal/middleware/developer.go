package middleware

import (
	"fmt"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type developerAuth struct {
	roleAuther
	client *repo.Oauth2ClientRepository
}

func newDeveloperAuth() roleAuther {
	return &developerAuth{
		client: repo.NewOauth2ClientRepository(),
	}
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
	clientID := ctx.Get("clientID").(string)
	client, err := auth.client.GetByUUID(uuid.MustParse(clientID))
	if err != nil {
		return echo.ErrUnauthorized
	}
	if !client.Approved {
		return echo.ErrUnauthorized
	}
	if !auth.isPathAllowed(ctx.Request().URL.Path) {
		return echo.ErrUnauthorized
	}
	if err := validateRequiredHeaders(ctx); err != nil {
		return err
	}
	platform := ctx.Request().Header.Get(X_TN_PLATFORM)
	if _, ok := platformsAuth[Platform(platform)]; !ok {
		return fmt.Errorf("platform %s is not supported", platform)
	}
	if err := platformsAuth[Platform(platform)].authorize(ctx); err != nil {
		return err
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
