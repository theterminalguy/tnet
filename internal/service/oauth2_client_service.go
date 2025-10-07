package service

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/theterminalguy/tenlog"
	"github.com/theterminalguy/tnet/ent"
	"github.com/theterminalguy/tnet/ent/oauth2client"
	"github.com/theterminalguy/tnet/ent/schema/userrole"
	"github.com/theterminalguy/tnet/ent/user"
	repo "github.com/theterminalguy/tnet/internal/repository"
	"github.com/theterminalguy/tnet/util"
	"github.com/theterminalguy/tnet/util/photo"
	"golang.org/x/crypto/bcrypt"
)

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

type Oauth2ClientService struct {
	Oauth2Repo *repo.Oauth2ClientRepository
	UserRepo   *repo.UserRepository
}

func NewOauth2ClientService() *Oauth2ClientService {
	return &Oauth2ClientService{
		Oauth2Repo: repo.NewOauth2ClientRepository(),
		UserRepo:   repo.NewUserRepository(),
	}
}

func (o2 *Oauth2ClientService) RegisterClient(p Oauth2ClientRegistraionParams) (*Oauth2ClientRegistrationResponse, error) {
	if err := repo.ValidateParams(p); err != nil {
		return nil, err
	}
	// Check if the email is already in use
	user, err := o2.UserRepo.GetByEmail(p.Contact.Email)
	if user != nil {
		return nil, ErrEmailAlreadyInUse
	}
	if err != nil && !ent.IsNotFound(err) {
		return nil, util.LogAndReturnErr("userRepo.GetByEmail", err)
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
	sr := repo.NewOauth2ScopeRepository()
	for _, scope := range p.AppInfo.Scopes {
		if !sr.IsValid(scope) {
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
	secret, err := util.SecureRandomHex(32) // 256 bits
	if err != nil {
		return nil, err
	}
	// Hash the secret
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	p.AppInfo.HashedSecret = string(hashedSecret)
	app, err := o2.Oauth2Repo.Register(p.AppInfo, p.Contact)
	if err != nil {
		return nil, err
	}
	return &Oauth2ClientRegistrationResponse{
		ClientID:     app.ID.String(),
		ClientSecret: secret,
		Scopes:       app.Scopes,
		GrantTypes:   o2.Oauth2Repo.GrantTypes(app),
	}, nil
}

func (o2 *Oauth2ClientService) ApproveClient(c ent.Oauth2Client) error {
	if !c.Approved && c.DeletedAt == nil {
		err := o2.Oauth2Repo.UpdateFields(&c, map[string]interface{}{
			oauth2client.FieldApproved: true,
		})
		if err != nil {
			return err
		}
	}
	u, err := o2.Oauth2Repo.GetUser(&c)
	if err != nil {
		return err
	}
	if !u.Approved && u.DeletedAt == nil {
		o2.UserRepo.UpdateFields(u, map[string]interface{}{
			user.FieldApproved: true,
		})
	}
	return nil
}
