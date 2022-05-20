package globalctx

import "github.com/labstack/echo/v4"

const (
	Platform       = "platform"
	PlatformTeamID = "platformID"
	PlatformUserID = "platformUserID"
)

func SetPlatformContext(ctx echo.Context, platform string) {
	ctx.Set(Platform, platform)
}

func SetPlatformTeamIDContext(ctx echo.Context, teamID string) {
	ctx.Set(PlatformTeamID, teamID)
}

func SetPlatformUserIDContext(ctx echo.Context, userID string) {
	ctx.Set(PlatformUserID, userID)
}

func GetPlatform(ctx echo.Context) string {
	return ctx.Get(Platform).(string)
}

func GetPlatformTeamID(ctx echo.Context) string {
	return ctx.Get(PlatformTeamID).(string)
}

func GetPlatformUserID(ctx echo.Context) string {
	return ctx.Get(PlatformUserID).(string)
}
