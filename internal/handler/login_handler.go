package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/randutil"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Token string

const (
	IDToken     = "id_token"
	AccessToken = "access_token"
)

var stateToken string

var gconf *oauth2.Config = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	RedirectURL:  fmt.Sprintf("%s/oauth2/google/callback", os.Getenv("APP_HOST")),
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

func GoogleOauth2CallbackHandler(c echo.Context) error {
	// TODO: we should redirect to a page once the flow completes without errors
	// this page should show the user profile obtained from google and also
	// shows their access token. They should be able to copy the access token to a clipboard.
	// We can implement a simple webpage with a button to copy the access token to clipboard.
	if stateToken != c.QueryParam("state") {
		// the state token is invalid, someone may be trying to intercept our login flow
		return echo.ErrUnauthorized
	}
	code := c.QueryParam("code")
	// TODO: avoid using context.Background()
	tok, err := gconf.Exchange(context.Background(), code)
	if err != nil {
		return err
	}
	client := gconf.Client(context.Background(), tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var userInfo repo.UserParams
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return err
	}
	ur := repo.NewUserRepository()
	record, err := ur.GetByEmail(userInfo.Email)
	if err != nil {
		if ent.IsNotFound(err) {
			userInfo.Role = userrole.Talent
			record, err = ur.Create(userInfo)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	claim := struct {
		Role      string `json:"role"`
		TokenType Token  `json:"token_type"`
		jwt.StandardClaims
	}{
		Role:      string(record.Role),
		TokenType: IDToken,
		StandardClaims: jwt.StandardClaims{
			Audience:  c.Request().Host,
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    c.Request().Host,
			Subject:   record.UUID.String(),
		},
	}
	jwttok := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := jwttok.SignedString([]byte(os.Getenv("JWT_SIGNED_SECRET")))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func TalentLoginHanlder(c echo.Context) error {
	stateToken = randutil.GenerateOauthStateToken()
	url := gconf.AuthCodeURL(stateToken)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}
