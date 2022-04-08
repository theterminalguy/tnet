package handler

import (
	"fmt"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/internal/middleware"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetCurrentUser(c echo.Context) (*ent.User, error) {
	token, err := middleware.ExtractToken(c)
	if err != nil {
		return nil, err
	}
	claims := token.Claims.(jwt.MapClaims)
	sub := claims["sub"]
	userUUID := fmt.Sprintf("%v", sub)

	ur := repo.NewUserRepository()
	record, err := ur.GetByID(uuid.MustParse(userUUID))
	if err != nil {
		return nil, err
	}
	return record, nil
}
