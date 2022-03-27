package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/slackappinstall"
	"github.com/google/uuid"
)

type SlackAppInstallRepository struct{}

type SlackAppInstallParams struct {
	UserID              uuid.UUID
	TeamID              string `json:"team_id" validate:"required"`
	TeamName            string `json:"team_name" validate:"required"`
	AuthedUserID        string `json:"authed_user_id" validate:"required"`
	AuthedUserEmail     string `json:"email" validate:"required,email"`
	AuthedUserTitle     string `json:"title"`
	AuthedUserPhone     string `json:"phone"`
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

func (*SlackAppInstallRepository) GetByID(id uuid.UUID) (*ent.SlackAppInstall, error) {
	record, err := dBConn.SlackAppInstall.Query().
		Where(slackappinstall.ID(id)).
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
		SetAuthedUserTitle(p.AuthedUserTitle).
		SetAuthedUserPhone(p.AuthedUserPhone).
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
	record, err := r.GetByID(id)
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

func (r *SlackAppInstallRepository) DeleteByID(id uuid.UUID) error {
	record, err := r.GetByID(id)
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
