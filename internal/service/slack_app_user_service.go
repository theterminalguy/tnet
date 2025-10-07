package service

import (
	"github.com/theterminalguy/tentn/ent"
	repo "github.com/theterminalguy/tentn/internal/repository"
	"github.com/google/uuid"
	"github.com/slack-go/slack"
)

type SlackAppUserService struct {
	AppUser    *repo.SlackAppUserRepository
	AppInstall *repo.SlackAppInstallRepository
}

func NewSlackAppUserService() *SlackAppUserService {
	return &SlackAppUserService{
		AppUser:    repo.NewSlackAppUserRepository(),
		AppInstall: repo.NewSlackAppInstallRepository(),
	}
}

func (s *SlackAppUserService) GetUserInfo(slackTeamID, slackUserID string) (*slack.User, error) {
	install, err := s.AppInstall.GetByTeamID(slackTeamID)
	if err != nil {
		return nil, err
	}
	api := slack.New(install.AccessToken)
	user, err := api.GetUserInfo(slackUserID)
	if err != nil {
		return nil, err
	}
	if user.Profile.Email == "" {
		userProfile, err := api.GetUserProfile(&slack.GetUserProfileParameters{
			UserID: slackUserID,
		})
		if err != nil {
			return nil, err
		}
		user.Profile.Email = userProfile.Email
	}
	return user, nil
}

func (s *SlackAppUserService) CreateUser(slackTeamID, slackUserID string, installID uuid.UUID) (*ent.SlackAppUser, error) {
	/**
	This is like a log that keeps track of every user in a Slack Workspace
	Using the app.

	If the User is found, we return the user
	else we fetch that user from slack and create a record in our DB
	*/
	user, err := s.AppUser.GetBySlackUserID(slackTeamID, slackUserID)
	if user != nil {
		return user, nil
	}
	if ent.IsNotFound(err) {
		user, err := s.GetUserInfo(slackTeamID, slackUserID)
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
