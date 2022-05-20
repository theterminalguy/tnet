package header

import (
	"errors"

	"github.com/labstack/echo/v4"
)

const (
	X_TN_PLATFORM      = "X-TN-Platform"
	X_TN_SLACK_TEAM_ID = "X-TN-Slack-Team-ID"
	X_TN_SLACK_USER_ID = "X-TN-Slack-User-ID"
	X_TN_TALENT_ID     = "X-TN-Talent-ID"
	X_TN_RECRUITER_ID  = "X-TN-Recruiter-ID"

	X_TN_INTERNAL_USER_ID = "X-TN-Internal-User-ID"
	X_TN_INTERNAL_API_KEY = "X-TN-Internal-API-Key"
)

// ValidateHeaders checks if the request has the required headers
func ValidateRequiredHeaders(ctx echo.Context) error {
	if ctx.Request().Header.Get(X_TN_PLATFORM) == "" {
		return errors.New("platform missing")
	}
	return nil
}
