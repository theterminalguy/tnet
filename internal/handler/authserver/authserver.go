package authserver

import (
	"net/http"

	"github.com/10hourlabs/tentn/internal/service"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/labstack/echo/v4"
	"github.com/ory/fosite"
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

type tokenHandler struct {
	oauth fosite.OAuth2Provider
}

func Oauth2ClientTokenHandler(c echo.Context) error {
	manager := manage.NewDefaultManager()
}
