package platform

import (
	"github.com/labstack/echo/v4"
)

type Platform string

const (
	PlatformSlack = Platform("platform/slack")
	PlatformWeb   = Platform("platform/web")
)

type PlatformAuth interface {
	Authorize(ctx echo.Context) error
}

var Auth = map[Platform]PlatformAuth{
	PlatformSlack: NewSlackPlatformAuth(),
	PlatformWeb:   NewWebPlatformAuth(),
}
