package globalctx

import (
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func SetCurrentUserContext(ctx echo.Context) error {
	// TODO: save in cache
	// First try to get the record from cache
	// If the record is not in cache,
	// then get the record from database
	// and store the record in cache
	//
	// Suggested interim cache library:
	// https://github.com/allegro/bigcache
	token := ctx.Get("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)
	ur := repo.NewUserRepository()
	user, err := ur.GetByID(uuid.MustParse(userID))
	if err != nil {
		return err
	}
	ctx.Set("currentUser", user)
	if user.Role == userrole.Developer {
		clientID := claims["aud"].(string)
		ctx.Set("clientID", clientID)
	}
	return nil
}
