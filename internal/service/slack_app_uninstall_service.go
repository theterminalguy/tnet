package service

import (
	"errors"

	"github.com/theterminalguy/tentn/internal/repository"
)

type SlackAppUninstallParams struct {
	Event struct {
		Type string `json:"type"`
	}
	EventTime int64 `json:"event_time"`

	BotToken string `json:"token"`

	TeamID string `json:"team_id"`
	APPID  string `json:"api_app_id"`
}

type SlackAppUnInstallService struct {
	SlackInstallRepo *repository.SlackAppInstallRepository
}

func NewSlackAppUnInstallService() *SlackAppUnInstallService {
	return &SlackAppUnInstallService{}
}

//record.TeamID = fmt.Sprintf("deleted_%s_%d", record.TeamID, time.Now().Unix())
//record.DeletedAt = time.Unix(params.EventTime, 0)
// the next steps should happen in a transaction

// first, we want to delete all the users of the app in the `slack_app_users` table
/**
DELETE FROM slack_app_users
WHERE slack_team_id = <slack_team_id>
*/

// then we want to delete all the search queries made by everyone in the Slack Team in the `search_logs` table
/**
DELETE FROM search_logs
WHERE platform_team_id = <slack_team_id>
AND platform = 'platform/slack'
*/

// then we should delete every talent's collection in the `talent_collections` table for the Slack Team
// this is simply deleting the collection belonging to the Primary User of the Slack Team
/**
DELETE FROM talent_collections
WHERE user_id = <app.user_id>
*/

// then we want to delete every single user in the Slack Team in the `users` table with the following query
/**
DELETE FROM users
WHERE role = 'recruiter'
AND email like '%@<slack_team_id>.slack.com'
*/

// finally, we want to delete the app from the slack_app_install table
func (s *SlackAppUnInstallService) UnInstall(params *SlackAppUninstallParams) error {
	app, err := s.SlackInstallRepo.GetByTeamID(params.TeamID)
	if err != nil {
		return err
	}
	if app.AppID != params.APPID {
		return errors.New("invalid app id, app id does not match")
	}
	if err := s.SlackInstallRepo.DeleteByID(app.ID); err != nil {
		return err
	}
	return nil
}
