package platform

import (
	"errors"

	"github.com/10hourlabs/tenlog"
	"github.com/10hourlabs/tentn/ent/schema/billing"
	"github.com/10hourlabs/tentn/internal/middleware/globalctx"
	"github.com/10hourlabs/tentn/internal/middleware/header"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/labstack/echo/v4"
)

type SlackPlatformAuth struct {
	PlatformAuth
	sar  *repo.SlackAppInstallRepository
	saus *service.SlackAppUserService
}

func NewSlackPlatformAuth() PlatformAuth {
	return &SlackPlatformAuth{
		sar:  repo.NewSlackAppInstallRepository(),
		saus: service.NewSlackAppUserService(),
	}
}

// authorize checks if the request is from a slack platform
func (s SlackPlatformAuth) Authorize(ctx echo.Context) error {
	if err := s.validateHeaders(ctx); err != nil {
		return err
	}
	slackTeamID := ctx.Request().Header.Get(header.X_TN_SLACK_TEAM_ID)
	slackUserID := ctx.Request().Header.Get(header.X_TN_SLACK_USER_ID)
	// TODO:
	// This should be cached
	// Suggested interim cache library:
	// https://github.com/allegro/bigcache
	app, err := s.sar.GetByTeamID(slackTeamID)
	if err != nil {
		return echo.ErrUnauthorized
	}
	if app.PaymentPlan != billing.Free {
		// TODO:
		// We only support free plan for now
		return echo.ErrUnauthorized
	}
	if app.AuthedUserID != slackUserID {
		// TODO: discuss this with the team
		// the app is not currently used by the user who installed it
		// so attach the current request to the user who installed it
		tenlog.Warn("App not in use by a primary user")
	}

	// track the current user making the request
	// TODO: let's cache this
	if _, err := s.saus.CreateUser(slackUserID, app.ID); err != nil {
		tenlog.Error("Failed to create slack app user", err)
		return echo.ErrUnauthorized
	}

	// For free plans, every request is tied to the primary user
	// The primary user is the user who installed the app
	primaryUser, err := s.sar.GetRecruiterByTeamID(slackTeamID)
	if err != nil {
		return err
	}
	return globalctx.SetCurrentRecruiterContext(ctx, primaryUser)
}

// validateHeaders checks if the request has the required headers for Slack
func (SlackPlatformAuth) validateHeaders(ctx echo.Context) error {
	if ctx.Request().Header.Get(header.X_TN_SLACK_TEAM_ID) == "" {
		return errors.New("slack team id missing")
	}
	if ctx.Request().Header.Get(header.X_TN_SLACK_USER_ID) == "" {
		return errors.New("slack user id missing")
	}
	return nil
}
