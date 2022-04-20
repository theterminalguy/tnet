package userauth

import (
	"fmt"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	"github.com/10hourlabs/tentn/internal/middleware/header"
	"github.com/10hourlabs/tentn/internal/middleware/platform"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type DeveloperAuth struct {
	RoleAuther
	client *repo.Oauth2ClientRepository
}

func NewDeveloperAuth() RoleAuther {
	return &DeveloperAuth{
		client: repo.NewOauth2ClientRepository(),
	}
}

func (auth DeveloperAuth) Authorize(u *ent.User, ctx echo.Context) error {
	if u.Role != userrole.Developer {
		return echo.ErrUnauthorized
	}
	if !u.Approved {
		return echo.ErrUnauthorized
	}
	// TODO:
	// First try to get the u from cache
	// If the u is not in cache,
	// then get the u from database
	// and store the u in cache
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
	if err := header.ValidateRequiredHeaders(ctx); err != nil {
		return err
	}
	p := ctx.Request().Header.Get(header.X_TN_PLATFORM)
	if _, ok := platform.Auth[platform.Platform(p)]; !ok {
		return fmt.Errorf("platform %s is not supported", p)
	}
	if err := platform.Auth[platform.Platform(p)].Authorize(ctx); err != nil {
		return err
	}
	// TODO: cater for internal/external clients
	// There are certain things internal clients can do that external clients can't
	return nil
}

func (DeveloperAuth) isPathAllowed(_ string) bool {
	// TODO: In the future,
	// there may be paths that are not allowed for developers

	// TODO: cater for internal/external clients
	// There are certain things internal clients can do that external clients can't
	return true
}
