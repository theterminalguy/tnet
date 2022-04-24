package authserver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/10hourlabs/tenlog"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/internal/tokgen"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/go-oauth2/oauth2/v4"
	oauth2_error "github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	manager *manage.Manager
	srv     *server.Server
)

func init() {
	manager = manage.NewDefaultManager()
	manager.SetClientTokenCfg(&manage.Config{
		IsGenerateRefresh: true,
		AccessTokenExp:    tokgen.DefaultAccessTokenExp,
		RefreshTokenExp:   tokgen.DefaultRefreshTokenExp,
	})
	manager.SetRefreshTokenCfg(&manage.RefreshingConfig{
		IsRemoveAccess:     true,
		IsRemoveRefreshing: true,
		IsResetRefreshTime: true,
	})

	// configure token storage
	manager.MustTokenStorage(repo.NewOauth2TokenRepository(), nil)
	// configure client storage
	manager.MapClientStorage(repo.NewOauth2ClientRepository())
	fmt.Println("DefaultSinging Method:", tokgen.GetDefaultSigningMethod())
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate(
		tokgen.DefaultSigningKeyID,
		tokgen.GetDefualtSigningKey(),
		tokgen.GetDefaultSigningMethod(),
	))
	srv = server.NewDefaultServer(manager)
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
	srv.SetAllowedResponseType(oauth2.Token)

	// Set the allowed grant types
	// https://oauth.net/2/grant-types
	//
	// We only the client credentials grant type and a refresh token grant type
	srv.SetAllowedGrantType(oauth2.ClientCredentials, oauth2.Refreshing)
	srv.SetClientScopeHandler(authorizeClientScope)
	srv.SetClientAuthorizedHandler(authorizeClientRequest)
}

// Oauth2ClienRepository allows for registration of Oauth2 client using
// standard Oauth2 and OpenID Connect flows.
func Oauth2RegistrationHandler(c echo.Context) error {
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

func Oauth2AccessTokenHandler(c echo.Context) error {
	// Parse the form request
	c.Request().ParseForm()
	return srv.HandleTokenRequest(c.Response(), c.Request())
}

func Oauth2RefreshTokenHandler(c echo.Context) error {
	// Parse the form request
	c.Request().ParseForm()
	// check for the refresh_token grant type
	reqGrantType := c.Request().Form.Get("grant_type")
	if reqGrantType != oauth2.Refreshing.String() {
		return c.String(http.StatusBadRequest, "Invalid grant_type")
	}
	return srv.HandleTokenRequest(c.Response(), c.Request())
}

func authorizeClientRequest(clientID string, grant oauth2.GrantType) (allowed bool, err error) {
	allowed = false
	err = oauth2_error.ErrUnsupportedGrantType
	if grant == oauth2.ClientCredentials || grant == oauth2.Refreshing {
		allowed = true
		err = nil
	}
	return
}

func authorizeClientScope(tgr *oauth2.TokenGenerateRequest) (bool, error) {
	client, err := repo.NewOauth2ClientRepository().GetByUUID(uuid.MustParse(tgr.ClientID))
	if err != nil {
		return false, err
	}
	tgr.UserID = client.UserID.String()
	// tgr.Scope is the requested scope
	if tgr.Scope != "" {
		reqScopes := strings.Split(tgr.Scope, ",")

		// Check if the requested scope is allowed
		for _, reqScope := range reqScopes {
			if !collection.Contains(client.Scopes, reqScope) {
				tenlog.Debug("Requested scope not allowed:", reqScope)
				return false, fmt.Errorf("requested scope %s not in client scope", reqScope)
			}
		}
	} else {
		tgr.Scope = strings.Join(client.Scopes, ",")
	}
	return true, nil
}
