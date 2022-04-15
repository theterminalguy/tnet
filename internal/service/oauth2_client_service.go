package service

import repo "github.com/10hourlabs/tentn/internal/repository"

// TODO: validate redirect_uris
// TODO: validate scopes
// TODO: validate client_type
// TODO: validate uris (app_logo_uri, app_homepage_uri)
// TODO: generate hashed_secret
// TODO: set grant_types
// If the client is a confidential client, set the grant_types to "authorization_code", "refresh_token", and "client_credentials".
// If the client is a public client, set the grant_types to "implicit"
type Oauth2ClientRegistraionParams struct {
	Contact repo.UserParams         `json:"contact"`
	AppInfo repo.Oauth2ClientParams `json:"app_info"`
}

type Oauth2ClientService struct{}

func NewOauth2ClientService() *Oauth2ClientService {
	return &Oauth2ClientService{}
}

func (*Oauth2ClientService) RegisterClient(p Oauth2ClientRegistraionParams) error {
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
