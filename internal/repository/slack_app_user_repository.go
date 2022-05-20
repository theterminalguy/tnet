package repository

import (
	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/slackappuser"
	"github.com/google/uuid"
)

type SlackAppUserRepository struct{}

type SlackAppUserParams struct {
	SlackAppInstallID uuid.UUID
	FullName          string `json:"full_name" validate:"required"`
	Title             string `json:"title"`
	Email             string `json:"email" validate:"required"`
	PhotoURL          string `json:"photo_url" validate:"required"`
	SlackUserID       string `json:"slack_user_id" validate:"required"`
	SlackTeamID       string `json:"slack_team_id" validate:"required"`
	Timezone          string `json:"timezone" validate:"required"`
	TimezoneLabel     string `json:"timezone_label" validate:"required"`
	IsBotUser         bool   `json:"is_bot_user"`
	Locale            string `json:"locale"`
}

func NewSlackAppUserRepository() *SlackAppUserRepository {
	return &SlackAppUserRepository{}
}

func (*SlackAppUserRepository) GetAll() ([]*ent.SlackAppUser, error) {
	records, err := dBConn.SlackAppUser.Query().
		Where(slackappuser.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*SlackAppUserRepository) GetByID(id uuid.UUID) (*ent.SlackAppUser, error) {
	record, err := dBConn.SlackAppUser.Query().
		Where(slackappuser.ID(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*SlackAppUserRepository) GetBySlackUserID(id string) (*ent.SlackAppUser, error) {
	record, err := dBConn.SlackAppUser.Query().
		Where(slackappuser.SlackUserID(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*SlackAppUserRepository) Create(p SlackAppUserParams) (*ent.SlackAppUser, error) {
	err := ValidateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := dBConn.SlackAppUser.
		Create().
		SetSlackAppInstallID(p.SlackAppInstallID).
		SetFullName(p.FullName).
		SetTitle(p.Title).
		SetEmail(p.Email).
		SetPhotoURL(p.PhotoURL).
		SetSlackUserID(p.SlackUserID).
		SetSlackTeamID(p.SlackTeamID).
		SetTimezone(p.Timezone).
		SetTimezoneLabel(p.TimezoneLabel).
		SetIsBotUser(p.IsBotUser).
		SetLocale(p.Locale).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *SlackAppUserRepository) Update(id uuid.UUID, p SlackAppUserParams) (*ent.SlackAppUser, error) {
	return nil, nil
}

func (r *SlackAppUserRepository) DeleteByUUID(id uuid.UUID) error {
	return nil
}
