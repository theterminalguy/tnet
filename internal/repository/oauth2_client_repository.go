package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/oauth2client"
	"github.com/google/uuid"
	"github.com/ory/fosite"
)

type Oauth2ClientRepository struct{}

type Oauth2ClientParams struct {
	AppName             string   `json:"app_name" validate:"required"`
	AppDescription      string   `json:"app_description" validate:"required"`
	AppLogoURI          string   `json:"app_logo_uri" validate:"required,url"`
	AppHomepageURI      string   `json:"app_homepage_uri" validate:"required,url"`
	AppPrivacyPolicyURI string   `json:"app_privacy_policy_uri" validate:"required,url"`
	ClientType          string   `json:"client_type" validate:"required"`
	Scopes              []string `json:"scopes" validate:"required"`
	RedirectURIs        []string `json:"redirect_uris" validate:"required"`

	UserID       uuid.UUID
	HashedSecret string
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

func (*Oauth2ClientRepository) GetByID(id uuid.UUID) (*ent.Oauth2Client, error) {
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
		SetID(uuid.New()).
		SetUserID(p.UserID).
		SetHashedSecret(p.HashedSecret).
		SetAppName(p.AppName).
		SetAppDescription(p.AppDescription).
		SetAppLogoURI(p.AppLogoURI).
		SetAppHomepageURI(p.AppHomepageURI).
		SetAppPrivacyPolicyURI(p.AppPrivacyPolicyURI).
		SetClientType(oauth2client.ClientType(p.ClientType)).
		SetScopes(p.Scopes).
		SetRedirectUris(p.RedirectURIs).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (*Oauth2ClientRepository) Register(c Oauth2ClientParams, u UserParams) (*ent.Oauth2Client, error) {
	tx, err := dBConn.Tx(dBContext)
	if err != nil {
		return nil, fmt.Errorf("starting a transaction: %w", err)
	}
	dev, err := tx.User.Create().
		SetFirstName(u.FirstName).
		SetLastName(u.LastName).
		SetPhotoURL(u.PhotoURL).
		SetEmail(u.Email).
		SetRole(u.Role).
		Save(dBContext)
	if err != nil {
		return nil, rollback(tx, fmt.Errorf("failed creating the user: %w", err))
	}
	o2, err := tx.Oauth2Client.
		Create().
		SetID(uuid.New()).
		SetUserID(dev.ID).
		SetHashedSecret(c.HashedSecret).
		SetAppName(c.AppName).
		SetAppDescription(c.AppDescription).
		SetAppLogoURI(c.AppLogoURI).
		SetAppHomepageURI(c.AppHomepageURI).
		SetAppPrivacyPolicyURI(c.AppPrivacyPolicyURI).
		SetClientType(oauth2client.ClientType(c.ClientType)).
		SetScopes(c.Scopes).
		SetRedirectUris(c.RedirectURIs).
		Save(dBContext)
	if err != nil {
		return nil, rollback(tx, fmt.Errorf("failed creating the oauth2 client: %w", err))
	}
	return o2, tx.Commit()
}

func (r *Oauth2ClientRepository) Update(id uuid.UUID, p Oauth2ClientParams) (*ent.Oauth2Client, error) {
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

func (r *Oauth2ClientRepository) DeleteByUUID(id uuid.UUID) error {
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

func (*Oauth2ClientRepository) GrantTypes(record *ent.Oauth2Client) []string {
	if record.ClientType == oauth2client.ClientTypePublic {
		return []string{"implicit"}
	}
	return []string{
		"authorization_code",
		"refresh_token",
		"client_credentials",
	}
}

func (*Oauth2ClientRepository) GetResponseTypes(record *ent.Oauth2Client) []string {
	if record.ClientType == oauth2client.ClientTypePublic {
		return []string{"code"}
	}
	return []string{
		"code",
		"token",
	}
}

// implement the Oauth fosite Client Manager interface

func (cr *Oauth2ClientRepository) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	c, err := cr.GetByID(uuid.MustParse(id))
	if err != nil {
		return nil, err
	}
	return Oauth2FositeClient{
		GrantTypes:    cr.GrantTypes(c),
		ResponseTypes: cr.GetResponseTypes(c),
		Oauth2Client:  c,
	}, nil
}

func (cr *Oauth2ClientRepository) ClientAssertionJWTValid(ctx context.Context, jti string) error {
	return nil
}

func (cr *Oauth2ClientRepository) SetClientAssertionJWT(ctx context.Context, jti string, exp time.Time) error {
	return nil
}
