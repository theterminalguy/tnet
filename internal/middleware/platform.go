package middleware

import (
	"github.com/labstack/echo/v4"
)

type Platform string

const (
	PlatformSlack = Platform("platform/slack")
	PlatformWeb   = Platform("platform/web")
)

type platformAuth interface {
	authorize(ctx echo.Context) error
}

var platformsAuth = map[Platform]platformAuth{
	PlatformSlack: newSlackPlatformAuth(),
	PlatformWeb:   newWebPlatformAuth(),
}
