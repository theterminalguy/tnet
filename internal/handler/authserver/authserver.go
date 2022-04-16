package authserver

import (
	"net/http"

	"github.com/10hourlabs/tentn/internal/service"
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
	//manager := manage.NewDefaultManager()
}
