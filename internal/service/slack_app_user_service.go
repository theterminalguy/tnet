package service

import (
	"os"

	"github.com/10hourlabs/tentn/ent"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/google/uuid"
	"github.com/slack-go/slack"
)

type SlackAppUserService struct {
	AppUser *repo.SlackAppUserRepository
}

func NewSlackAppUserService() *SlackAppUserService {
	return &SlackAppUserService{
		AppUser: repo.NewSlackAppUserRepository(),
	}
}

func (*SlackAppUserService) GetUserInfo(slackUserID string) (*slack.User, error) {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	return api.GetUserInfo(slackUserID)
}

func (s *SlackAppUserService) CreateUser(slackUserID string, installID uuid.UUID) (*ent.SlackAppUser, error) {
	user, err := s.AppUser.GetBySlackUserID(slackUserID)
	if user != nil {
		return user, nil
	}
	if ent.IsNotFound(err) {
		user, err := s.GetUserInfo(slackUserID)
		if err != nil {
			return nil, err
		}
		params := repo.SlackAppUserParams{
			SlackAppInstallID: installID,
			FullName:          user.RealName,
			Title:             user.Profile.Title,
			Email:             user.Profile.Email,
			PhotoURL:          user.Profile.Image72,
			SlackUserID:       user.ID,
			SlackTeamID:       user.TeamID,
			Timezone:          user.TZ,
			TimezoneLabel:     user.TZLabel,
			IsBotUser:         user.IsBot,
			Locale:            user.Locale,
		}
		return s.AppUser.Create(params)
	}
	return nil, err
}
