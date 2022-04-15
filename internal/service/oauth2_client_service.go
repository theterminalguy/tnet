package service

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/10hourlabs/tenlog"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/util/photo"
	"github.com/ory/fosite"
)

// TODO: validate scopes
// TODO: validate client_type
// TODO: generate hashed_secret
// TODO: set grant_types
// If the client is a confidential client, set the grant_types to "authorization_code", "refresh_token", and "client_credentials".
// If the client is a public client, set the grant_types to "implicit"

var ErrEmailAlreadyInUse = errors.New("the provided email is already in use")
var ErrInvalidRedirectURI = errors.New("invalid redirect_uri")

type Oauth2ClientRegistraionParams struct {
	Contact repo.UserParams         `json:"contact" validate:"required"`
	AppInfo repo.Oauth2ClientParams `json:"app_info" validate:"required"`
}

type Oauth2ClientService struct{}

func NewOauth2ClientService() *Oauth2ClientService {
	return &Oauth2ClientService{}
}

func (*Oauth2ClientService) RegisterClient(p Oauth2ClientRegistraionParams) error {
	userRepo := repo.NewUserRepository()

	// Check if the email is already in use
	user, err := userRepo.GetByEmail(p.Contact.Email)
	if user != nil {
		return ErrEmailAlreadyInUse
	}
	tenlog.Error("userRepo.GetByEmail", err)

	// validate redirect_uris
	var errors []error
	uris := p.AppInfo.RedirectURIs
	for _, uri := range uris {
		rURI, err := url.Parse(uri)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		if !fosite.IsRedirectURISecureStrict(rURI) {
			errors = append(errors, ErrInvalidRedirectURI)
		}
	}
	if len(errors) > 0 {
		tenlog.Error("invalid redirect_uris", errors)
		return fmt.Errorf("%v", errors)
	}

	// validate scopes
	p.Contact.Role = userrole.Developer
	p.Contact.PhotoURL = photo.GenerateDefaultPhoto(p.Contact.FirstName, p.Contact.LastName)

	/*err := ValidateParams(p)
	if err != nil {
		return nil, err
	}
	h := repo.NewOauth2ClientRepository()
	record, err := h.Create(p.AppInfo)
	if err != nil {
		return nil, err
	}
	return record, nil*/
	return nil
}
