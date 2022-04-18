package handler

import (
	"fmt"

	"github.com/10hourlabs/tentn/ent"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetCurrentUser(c echo.Context) (*ent.User, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	sub := claims["sub"]
	userID := fmt.Sprintf("%v", sub)

	ur := repo.NewUserRepository()
	record, err := ur.GetByID(uuid.MustParse(userID))
	if err != nil {
		return nil, err
	}
	return record, nil
}
