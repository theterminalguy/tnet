package platform

import (
	"errors"

	"github.com/10hourlabs/tentn/internal/middleware/globalctx"
	"github.com/10hourlabs/tentn/internal/middleware/header"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type WebPlatformAuth struct {
	PlatformAuth
	ur *repo.UserRepository
}

func NewWebPlatformAuth() PlatformAuth {
	return &WebPlatformAuth{
		ur: repo.NewUserRepository(),
	}
}

// authorize checks if the request is from a web platform
func (w WebPlatformAuth) Authorize(ctx echo.Context) error {
	if err := w.validateHeaders(ctx); err != nil {
		return err
	}
	talentID := ctx.Request().Header.Get(header.X_TN_TALENT_ID)
	if err := globalctx.SetCurrentTalentContext(ctx, uuid.MustParse(talentID)); err != nil {
		return err
	}
	recruiterID := ctx.Request().Header.Get(header.X_TN_RECRUITER_ID)
	recruiter, err := w.ur.GetRecruiterByID(uuid.MustParse(recruiterID))
	if err != nil {
		return err
	}
	globalctx.SetPlatformContext(ctx, string(PlatformWeb))
	globalctx.SetPlatformTeamIDContext(ctx, recruiter.ID.String())
	globalctx.SetPlatformUserIDContext(ctx, recruiter.ID.String())

	return globalctx.SetCurrentRecruiterContext(ctx, recruiter)
}

// validateHeaders checks if the request has the required headers for Web
func (WebPlatformAuth) validateHeaders(ctx echo.Context) error {
	headers := ctx.Request().Header
	if headers.Get(header.X_TN_TALENT_ID) == "" {
		return errors.New("talent id missing")
	}
	if headers.Get(header.X_TN_RECRUITER_ID) == "" {
		return errors.New("recruiter id missing")
	}
	return nil
}
