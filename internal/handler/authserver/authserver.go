package authserver

import (
	"net/http"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
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
	manager := manage.NewDefaultManager()
	// configure token storage
	manager.MustTokenStorage(repo.NewOauth2TokenRepository(), nil)
	// configure client storage
	manager.MapClientStorage(repo.NewOauth2ClientRepository())
	srv := server.NewDefaultServer(manager)

	// Do not allow GET request for token
	srv.SetAllowGetAccessRequest(false)

	// Set the authorization code expiration time.
	srv.SetClientInfoHandler(server.ClientFormHandler)

}
