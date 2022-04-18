package authserver

import (
	"net/http"

	"github.com/10hourlabs/tenlog"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/go-oauth2/oauth2/v4"
	oauth2_error "github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/labstack/echo/v4"
)

// Oauth2ClienRepository allows for registration of Oauth2 client using
// standard Oauth2 and OpenID Connect flows.
func Oauth2ClientRegisterationHandler(c echo.Context) error {
	params := new(service.Oauth2ClientRegistraionParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	oauth2Service := service.NewOauth2ClientService()
	resp, err := oauth2Service.RegisterClient(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, resp)
}

func Oauth2ClientTokenHandler(c echo.Context) error {
	// Parse the form request
	c.Request().ParseForm()

	manager := manage.NewDefaultManager()
	// configure token storage
	manager.MustTokenStorage(repo.NewOauth2TokenRepository(), nil)
	// configure client storage
	manager.MapClientStorage(repo.NewOauth2ClientRepository())
	srv := server.NewDefaultServer(manager)

	srv.SetInternalErrorHandler(func(err error) (re *oauth2_error.Response) {
		tenlog.Debug("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *oauth2_error.Response) {
		tenlog.Debug("Response Error:", re)
	})

	// Do not allow GET request for token
	srv.SetAllowGetAccessRequest(false)

	// This default handler let's us get the client_id and client_secret from the request
	srv.SetClientInfoHandler(server.ClientFormHandler)

	// Set the allowed response types
	// Only allow the token response type
	srv.SetAllowedResponseType([]oauth2.ResponseType{oauth2.Token}...)

	// Set the allowed grant types
	// https://oauth.net/2/grant-types
	//
	// We only the client credentials grant type
	srv.SetAllowedGrantType([]oauth2.GrantType{oauth2.ClientCredentials}...)

	srv.SetClientAuthorizedHandler(authorizeClientRequest)
	return srv.HandleTokenRequest(c.Response(), c.Request())
}

func authorizeClientRequest(clientID string, grant oauth2.GrantType) (allowed bool, err error) {
	allowed = false
	err = oauth2_error.ErrUnsupportedGrantType
	if grant == oauth2.ClientCredentials {
		allowed = true
		err = nil
	}
	return
}
