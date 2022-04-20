package platform

import (
	"errors"

	"github.com/10hourlabs/tentn/internal/middleware/header"
	"github.com/labstack/echo/v4"
)

type WebPlatformAuth struct {
	PlatformAuth
}

func NewWebPlatformAuth() PlatformAuth {
	return &WebPlatformAuth{}
}

// authorize checks if the request is from a web platform
func (w WebPlatformAuth) Authorize(ctx echo.Context) error {
	if err := w.validateHeaders(ctx); err != nil {
		return err
	}
	return nil
}

// validateHeaders checks if the request has the required headers for Web
func (WebPlatformAuth) validateHeaders(ctx echo.Context) error {
	headers := ctx.Request().Header
	if headers.Get(header.X_TN_WEB_USER_ID) == "" {
		return errors.New("web user id missing")
	}
	return nil
}
