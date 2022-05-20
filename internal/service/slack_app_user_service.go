package service

import (
	"os"

	"github.com/slack-go/slack"
)

type SlackAppUserService struct{}

func NewSlackAppUserService() *SlackAppUserService {
	return &SlackAppUserService{}
}

func (s *SlackAppUserService) GetUserInfo(slackUserID string) (*slack.User, error) {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	return api.GetUserInfo(slackUserID)
}
