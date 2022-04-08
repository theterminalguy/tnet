package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/10hourlabs/tentn/ent/schema/userrole"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func EnforceApprovedRecruiter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tok, err := ExtractToken(c)
			if err != nil {
				return err
			}
			claims := tok.Claims.(jwt.MapClaims)
			if claims["role"] != string(userrole.Recruiter) {
				return echo.ErrUnauthorized
			}
			// TODO: I noticed that after I manually approved a recruiter, I was still getting an unauthorized error.
			// This is because the previous jwt token had a approved field set to false.
			// So users will have to obtain a new token after they have been approved.
			// We fix this by removing the approved field from the jwt token and making a database call to get the approved field.
			if claims["approved"] == false {
				return echo.ErrUnauthorized
			}
			// TODO: prevent deleted accounts
			// This isn't done now as we are yet to figure out the best way to approach this
			// The obvious straightforward way is to check if the user is deleted, but that
			// would require a database query. We may have to cache the user data to allow for
			// faster access.
			return next(c)
		}
	}
}

func ExtractToken(c echo.Context) (*jwt.Token, error) {
	auth := c.Request().Header.Get("Authorization")
	if auth == "" {
		return nil, echo.ErrUnauthorized
	}
	tok := strings.SplitAfter(auth, "Bearer ")
	// parse token
	token, err := jwt.Parse(tok[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SIGNED_SECRET")), nil
	})
	return token, err
}
