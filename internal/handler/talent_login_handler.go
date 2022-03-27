package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/randutil"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauth2StateToken string

type GoogleOauth2Client struct {
	*oauth2.Config
}

type UserInfo struct {
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
	Email      string `json:"email"`
}

var gconf *GoogleOauth2Client = &GoogleOauth2Client{
	Config: &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		RedirectURL:  fmt.Sprintf("%s/oauth2/google/callback", os.Getenv("APP_HOST")),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	},
}

func (s *GoogleOauth2Client) GetUsersProfile(tok *oauth2.Token) (*repo.UserParams, error) {
	client := gconf.Client(context.Background(), tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return &repo.UserParams{
		FirstName: userInfo.GivenName,
		LastName:  userInfo.FamilyName,
		Email:     userInfo.Email,
		PhotoURL:  userInfo.Picture,
	}, nil
}

// GoogleOauth2CallbackHandler handles the callback from Google OAuth2
// anyone who signs up with Google will be redirected to this handler
// and will be automatically created with the "Talent" role
func GoogleOauth2CallbackHandler(c echo.Context) error {
	// TODO: we should redirect to a page once the flow completes without errors
	// this page should show the user profile obtained from google and also
	// shows their access token. They should be able to copy the access token to a clipboard.
	// We can implement a simple webpage with a button to copy the access token to clipboard.
	if googleOauth2StateToken != c.QueryParam("state") {
		// the state token is invalid, someone may be trying to intercept our login flow
		return echo.ErrUnauthorized
	}
	code := c.QueryParam("code")
	// TODO: avoid using context.Background()
	tok, err := gconf.Exchange(context.Background(), code)
	if err != nil {
		return err
	}
	talentProfile, err := gconf.GetUsersProfile(tok)
	if err != nil {
		return err
	}
	ts := service.NewTalentService()
	talent, err := ts.RegisterTalent(talentProfile)
	if err != nil {
		return err
	}
	claims := NewJWTClaims(talent, c)
	token, err := claims.GenerateToken()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func TalentLoginHanlder(c echo.Context) error {
	googleOauth2StateToken = randutil.GenerateOauthStateToken()
	url := gconf.AuthCodeURL(googleOauth2StateToken)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}
