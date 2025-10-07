package globalctx

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	repo "github.com/theterminalguy/tnet/internal/repository"
	"github.com/theterminalguy/tnet/internal/repository/scope"
)

func SetCurrentTalentContext(ctx echo.Context, userID uuid.UUID) error {
	// TODO: save in cache
	// First try to get the record from cache
	// If the record is not in cache,
	// then get the record from database
	// and store the record in cache
	//
	// Suggested interim cache library:
	// https://github.com/allegro/bigcache
	tr := repo.NewTalentRepository()
	talent, err := tr.GetTalentByUserID(userID)
	if err != nil {
		return err
	}
	currentTalent := scope.NewTalentScope(talent.GetTalent())
	ctx.Set("currentTalent", currentTalent)
	return nil
}
