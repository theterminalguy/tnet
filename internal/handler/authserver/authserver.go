package authserver

import (
	"log"
	"net/http"
	"time"

	"github.com/10hourlabs/tentn/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/handler/oauth2"
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
	// check the api docs of compose.Config for further configuration options
	//var oauth2 = compose.ComposeAllEnabled(config, store, secret, privateKey)
	var config = &compose.Config{
		AccessTokenLifespan: time.Minute * 30,
	}
	var oauth2Provider = compose.Compose(
		config,
		storage,
		compose.NewOAuth2JWTStrategy(
			key,
			// HMACStrategy is used to sign refresh token
			// therefore not required for our example
			nil,
		),
		// BCrypt hasher is automatically created when omitted.
		// Hasher is used to store hashed client authentication passwords.
		nil,
		compose.OAuth2ClientCredentialsGrantFactory,
	)
	ctx := c.Request().Context()
	session := &oauth2.JWTSession{}
	th := &tokenHandler{}
	ar, err := th.oauth.NewAccessRequest(ctx, c.Request(), session)
	if err != nil {
		// TODO: echo.Logger.Error(err)
		log.Println("An Error Occured: ", err)
		return err
	}
	resp, err := th.oauth.NewAccessResponse(ctx, ar)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}
