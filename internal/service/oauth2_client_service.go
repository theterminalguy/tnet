package service

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/10hourlabs/tenlog"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/util"
	"github.com/10hourlabs/tentn/util/photo"
	"golang.org/x/crypto/bcrypt"
)

// defaultBcryptWorkFactor is the default bcrypt work factor.
const defaultBCryptWorkFactor = 12

var ErrEmailAlreadyInUse = errors.New("the provided email is already in use")
var ErrInvalidRedirectURI = errors.New("invalid redirect_uri")

type Oauth2ClientRegistraionParams struct {
	Contact repo.UserParams         `json:"contact" validate:"required"`
	AppInfo repo.Oauth2ClientParams `json:"app_info" validate:"required"`
}

// Oauth2ClientRegistrationResponse is the response to the Oauth2ClientRegistrationRequest
// https://www.oauth.com/oauth2-servers/client-registration/client-id-secret/
type Oauth2ClientRegistrationResponse struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scopes       []string `json:"scopes"`
	GrantTypes   []string `json:"grant_types"`
}

type Oauth2ClientService struct{}

func NewOauth2ClientService() *Oauth2ClientService {
	return &Oauth2ClientService{}
}

func (*Oauth2ClientService) RegisterClient(p Oauth2ClientRegistraionParams) (*Oauth2ClientRegistrationResponse, error) {
	if err := repo.ValidateParams(p); err != nil {
		return nil, err
	}
	userRepo := repo.NewUserRepository()
	// Check if the email is already in use
	user, err := userRepo.GetByEmail(p.Contact.Email)
	if user != nil {
		return nil, ErrEmailAlreadyInUse
	}
	if err != nil {
		util.LogAndReturnErr("userRepo.GetByEmail", err)
	}
	// validate redirect_uris
	var errors []error
	uris := p.AppInfo.RedirectURIs
	for _, uri := range uris {
		rURI, err := url.Parse(uri)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		if !(util.IsValidRedirectURI(rURI) && util.IsRedirectURISecureStrict(rURI)) {
			errors = append(errors, ErrInvalidRedirectURI)
		}
	}
	if len(errors) > 0 {
		tenlog.Error("invalid redirect_uris", errors)
		return nil, fmt.Errorf("%v", errors)
	}
	// validate provided scopes
	for _, scope := range p.AppInfo.Scopes {
		if _, ok := repo.OauthScopes[scope]; !ok {
			errors = append(errors, fmt.Errorf("invalid scope: %s", scope))
		}
	}
	if len(errors) > 0 {
		tenlog.Error("invalid scopes", errors)
		return nil, fmt.Errorf("%v", errors)
	}
	p.Contact.Role = userrole.Developer
	p.Contact.PhotoURL = photo.GenerateDefaultPhoto(p.Contact.FirstName, p.Contact.LastName)
	// Genereate client secret
	secret, err := util.RandomBytes(32) // 256 bits
	if err != nil {
		return nil, err
	}
	// Hash the secret
	hashedSecret, err := bcrypt.GenerateFromPassword(secret, defaultBCryptWorkFactor)
	if err != nil {
		return nil, err
	}
	p.AppInfo.HashedSecret = string(hashedSecret)
	oauth2Repo := repo.NewOauth2ClientRepository()
	app, err := oauth2Repo.Register(p.AppInfo, p.Contact)
	if err != nil {
		return nil, err
	}
	// convert hex to string
	secretHex := fmt.Sprintf("%x", secret)
	return &Oauth2ClientRegistrationResponse{
		ClientID:     app.ID.String(),
		ClientSecret: secretHex,
		Scopes:       app.Scopes,
		GrantTypes:   oauth2Repo.GrantTypes(app),
	}, nil
}
