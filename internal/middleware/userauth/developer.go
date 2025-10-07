package userauth

import (
	"errors"
	"fmt"

	"github.com/theterminalguy/tenlog"
	"github.com/theterminalguy/tentn/ent"
	"github.com/theterminalguy/tentn/ent/schema/userrole"
	"github.com/theterminalguy/tentn/internal/middleware/header"
	"github.com/theterminalguy/tentn/internal/middleware/platform"
	repo "github.com/theterminalguy/tentn/internal/repository"
	"github.com/theterminalguy/tentn/util"
	"github.com/theterminalguy/tentn/util/collection"
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
		tenlog.Error("client not approved", "client", client.ID)
		return errors.New("client is not approved")
	}
	if err := auth.enforceOauth2Scope(client, ctx); err != nil {
		return err
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
	return nil
}

func (DeveloperAuth) enforceOauth2Scope(c *ent.Oauth2Client, ctx echo.Context) error {
	if c.IsInternal {
		return nil
	}
	path := ctx.Request().URL.Path
	verb := ctx.Request().Method
	reqScope := util.PathToOauth2Scope(path, verb)
	if reqScope != "" {
		if collection.Contains(c.Scopes, reqScope) {
			return nil
		}
	}
	return fmt.Errorf("client does not have the required scope %s", reqScope)
}
