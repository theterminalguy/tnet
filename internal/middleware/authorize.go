package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/theterminalguy/tenlog"
	"github.com/theterminalguy/tnet/ent"
	"github.com/theterminalguy/tnet/ent/schema/userrole"
	"github.com/theterminalguy/tnet/internal/middleware/globalctx"
	"github.com/theterminalguy/tnet/internal/middleware/userauth"
)

var roleAuth = map[userrole.Role]userauth.RoleAuther{
	userrole.Talent:    userauth.NewTalentAuth(),
	userrole.Recruiter: userauth.NewRecruiterAuth(),
	userrole.Developer: userauth.NewDeveloperAuth(),
}

func AuthorizieUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			if !token.Valid {
				return echo.ErrUnauthorized
			}
			c.Set("token", token)
			if err := globalctx.SetCurrentUserContext(c); err != nil {
				return echo.ErrUnauthorized
			}
			user := c.Get("currentUser").(*ent.User)
			if _, ok := roleAuth[user.Role]; !ok {
				tenlog.Error("role not found", "role", user.Role)
				return echo.ErrUnauthorized
			}
			auth := roleAuth[user.Role]
			if err := auth.Authorize(user, c); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			return next(c)
		}
	}
}
