package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/slackappinstall"
	"github.com/google/uuid"
)

type SlackAppInstallRepository struct{}

type SlackAppInstallParams struct {
	UserID              int
	TeamID              string `json:"team_id" validate:"required"`
	TeamName            string `json:"team_name" validate:"required"`
	AuthedUserID        string `json:"authed_user_id" validate:"required"`
	AuthedUserEmail     string `json:"email" validate:"required,email"`
	AppID               string `json:"app_id" validate:"required"`
	BotUserID           string `json:"bot_user_id" validate:"required"`
	AccessToken         string `json:"access_token" validate:"required"`
	TokenType           string `json:"token_type" validate:"required"`
	Scope               string `json:"scope" validate:"required"`
	IsEnterpriseInstall bool   `json:"is_enterprise_install"`
}

func NewSlackAppInstallRepository() *SlackAppInstallRepository {
	return &SlackAppInstallRepository{}
}

func (*SlackAppInstallRepository) GetAll() ([]*ent.SlackAppInstall, error) {
	records, err := dBConn.SlackAppInstall.Query().
		Where(slackappinstall.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*SlackAppInstallRepository) GetByUUID(id uuid.UUID) (*ent.SlackAppInstall, error) {
	record, err := dBConn.SlackAppInstall.Query().
		Where(slackappinstall.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*SlackAppInstallRepository) GetByEmail(email string) (*ent.SlackAppInstall, error) {
	record, err := dBConn.SlackAppInstall.Query().
		Where(slackappinstall.AuthedUserEmailEQ(email)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*SlackAppInstallRepository) Create(p SlackAppInstallParams) (*ent.SlackAppInstall, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := dBConn.SlackAppInstall.
		Create().
		SetUserID(p.UserID).
		SetTeamID(p.TeamID).
		SetTeamName(p.TeamName).
		SetAuthedUserID(p.AuthedUserID).
		SetAuthedUserEmail(p.AuthedUserEmail).
		SetAppID(p.AppID).
		SetBotUserID(p.BotUserID).
		SetAccessToken(p.AccessToken).
		SetTokenType(p.TokenType).
		SetScope(p.Scope).
		SetIsEnterpriseInstall(p.IsEnterpriseInstall).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *SlackAppInstallRepository) Update(id uuid.UUID, p SlackAppInstallParams) (*ent.SlackAppInstall, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := r.GetByUUID(id)
	if err != nil {
		return nil, err
	}
	_, err = record.Update().
		// TODO: set other fields here
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (r *SlackAppInstallRepository) DeleteByUUID(id uuid.UUID) error {
	record, err := r.GetByUUID(id)
	if err != nil {
		return err
	}
	_, err = record.Update().
		SetDeletedAt(time.Now()).
		Save(dBContext)
	if err != nil {
		return err
	}
	return nil
}
