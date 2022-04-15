package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/oauth2client"
	"github.com/google/uuid"
)

type Oauth2ClientRepository struct{}

type Oauth2ClientParams struct {
	DeveloperInfo  UserParams `json:"developer_info"`
	AppName        string     `json:"app_name"`
	AppDescription string     `json:"app_description"`
	AppLogoURI     string     `json:"app_logo_uri"`
	AppHomepageURI string     `json:"app_homepage_uri"`
	AppPrivacyURI  string     `json:"app_privacy_policy_uri"`
	ClientType     string     `json:"client_type"`
	Scopes         []string   `json:"scopes"`
}

func NewOauth2ClientRepository() *Oauth2ClientRepository {
	return &Oauth2ClientRepository{}
}

func (*Oauth2ClientRepository) GetAll() ([]*ent.Oauth2Client, error) {
	records, err := dBConn.Oauth2Client.Query().
		Where(oauth2client.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*Oauth2ClientRepository) GetByUUID(id uuid.UUID) (*ent.Oauth2Client, error) {
	record, err := dBConn.Oauth2Client.Query().
		Where(oauth2client.IDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*Oauth2ClientRepository) Create(p Oauth2ClientParams) (*ent.Oauth2Client, error) {
	err := ValidateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := dBConn.Oauth2Client.
		Create().
		// TODO: set other fields here
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *Oauth2ClientRepository) Update(id uuid.UUID, p Oauth2ClientParams) (*ent.Oauth2Client, error) {
	err := ValidateParams(p)
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

func (r *Oauth2ClientRepository) DeleteByUUID(id uuid.UUID) error {
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
