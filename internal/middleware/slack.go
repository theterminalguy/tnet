package middleware

import (
	"errors"

	"github.com/labstack/echo/v4"
)

type slackPlatformAuth struct {
	platformAuth
}

func newSlackPlatformAuth() platformAuth {
	return &slackPlatformAuth{}
}

// authorize checks if the request is from a slack platform
func (s slackPlatformAuth) authorize(ctx echo.Context) error {
	if err := s.validateHeaders(ctx); err != nil {
		return err
	}
	return nil
}

// validateHeaders checks if the request has the required headers for Slack
func (slackPlatformAuth) validateHeaders(ctx echo.Context) error {
	headers := ctx.Request().Header
	slackTeamID := headers.Get(X_TN_SLACK_TEAM_ID)
	slackUserID := headers.Get(X_TN_SLACK_USER_ID)
	if slackTeamID == "" {
		return errors.New("slack team id missing")
	}
	if slackUserID == "" {
		return errors.New("slack user id missing")
	}
	// TODO:
	// First try to get the user from cache
	// If the user is not in cache,
	// then get the user from database
	// and store the user in cache
	//
	// Suggested interim cache library:
	// https://github.com/allegro/bigcache
	return nil
}
