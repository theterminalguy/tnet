package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/slackappinstall"
	"github.com/google/uuid"
)

type SlackAppInstallRepository struct{}

type SlackAppInstallParams struct {
	UserID              uuid.UUID
	DeletedAt           *time.Time
	InstallCount        int64
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

func (*SlackAppInstallRepository) GetByTeamID(teamID string) (*ent.SlackAppInstall, error) {
	record, err := dBConn.SlackAppInstall.Query().
		Where(slackappinstall.TeamID(teamID)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*SlackAppInstallRepository) GetDeletedInstallation(teamID string) (*ent.SlackAppInstall, error) {
	record, err := dBConn.SlackAppInstall.Query().
		Where(slackappinstall.TeamID(teamID)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt == nil {
		return nil, errors.New("record not found")
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

func (*SlackAppInstallRepository) GetRecruiterByEmail(email string) (*ent.User, error) {
	record, err := dBConn.SlackAppInstall.Query().
		Where(slackappinstall.AuthedUserEmailEQ(email)).
		WithUser().
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record.Edges.User, nil
}

// GetRecruiterByTeamID returns the recruiter that installed the app
func (*SlackAppInstallRepository) GetRecruiterByTeamID(teamID string) (*ent.User, error) {
	record, err := dBConn.SlackAppInstall.Query().
		Where(slackappinstall.TeamID(teamID)).
		WithUser().
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record.Edges.User, nil
}

func (*SlackAppInstallRepository) Create(p SlackAppInstallParams) (*ent.SlackAppInstall, error) {
	err := ValidateParams(p)
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
	err := ValidateParams(p)
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

func (*SlackAppInstallRepository) UpdateFields(s *ent.SlackAppInstall, fields map[string]interface{}) error {
	bldr := s.Update()
	for k, v := range fields {
		switch k {
		case slackappinstall.FieldDeletedAt:
			bldr.ClearDeletedAt()
		case slackappinstall.FieldInstallCount:
			bldr.SetInstallCount(v.(int))
		case slackappinstall.FieldAccessToken:
			bldr.SetAccessToken(v.(string))
		case slackappinstall.FieldTeamName:
			bldr.SetTeamName(v.(string))
		case slackappinstall.FieldAuthedUserEmail:
			bldr.SetAuthedUserEmail(v.(string))
		case slackappinstall.FieldAuthedUserTitle:
			bldr.SetAuthedUserTitle(v.(string))
		case slackappinstall.FieldAuthedUserPhone:
			bldr.SetAuthedUserPhone(v.(string))
		case slackappinstall.FieldBotUserID:
			bldr.SetBotUserID(v.(string))
		case slackappinstall.FieldIsEnterpriseInstall:
			bldr.SetIsEnterpriseInstall(v.(bool))
		case slackappinstall.FieldScope:
			bldr.SetScope(v.(string))
		default:
			return fmt.Errorf("unknown field %s", k)
		}
	}
	if _, err := bldr.Save(dBContext); err != nil {
		return err
	}
	return nil
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
