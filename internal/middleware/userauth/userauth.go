package userauth

import (
	"github.com/labstack/echo/v4"
	"github.com/theterminalguy/tnet/ent"
)

type RoleAuther interface {
	Authorize(user *ent.User, c echo.Context) error
}
