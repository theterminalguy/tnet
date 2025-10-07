package userauth

import (
	"github.com/theterminalguy/tentn/ent"
	"github.com/labstack/echo/v4"
)

type RoleAuther interface {
	Authorize(user *ent.User, c echo.Context) error
}
