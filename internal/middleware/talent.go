package middleware

import (
	"strings"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/repository/scope"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
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
	return auth.setCurrentTalentContext(ctx, user.ID)
}

func (talentAuth) isPathAllowed(path string) bool {
	return strings.Contains(path, "/v1/talent")
}

func (talentAuth) setCurrentTalentContext(ctx echo.Context, userID uuid.UUID) error {
	// TODO:
	// First try to get the user from cache
	// If the user is not in cache,
	// then get the user from database
	// and store the user in cache
	//
	// Suggested interim cache library:
	// https://github.com/allegro/bigcache
	talent, err := repo.NewTalentRepository().GetTalentByUserID(userID)
	if err != nil {
		return err
	}
	currentTalent := scope.NewTalentScope(talent.GetTalent())
	ctx.Set(oneword.CurrentTalent, currentTalent)
	return nil
}
