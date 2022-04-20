package middleware

import (
	"errors"

	"github.com/labstack/echo/v4"
)

type webPlatformAuth struct {
	platformAuth
}

func newWebPlatformAuth() platformAuth {
	return &webPlatformAuth{}
}

// authorize checks if the request is from a web platform
func (w webPlatformAuth) authorize(ctx echo.Context) error {
	if err := w.validateHeaders(ctx); err != nil {
		return err
	}
	return nil
}

// validateHeaders checks if the request has the required headers for Web
func (webPlatformAuth) validateHeaders(ctx echo.Context) error {
	headers := ctx.Request().Header
	if headers.Get(X_TN_WEB_USER_ID) == "" {
		return errors.New("web user id missing")
	}
	return nil
}
