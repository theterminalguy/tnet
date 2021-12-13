package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/10hourlabs/tentn/randutil"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

var slackOauth2StateToken string

type SlackOauthResponse struct {
	AuthedUser struct {
		ID string `json:"id"`
	} `json:"authed_user"`
	Team struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"team"`
	Ok                  bool   `json:"ok"`
	AppID               string `json:"app_id"`
	Scope               string `json:"scope"`
	TokenType           string `json:"token_type"`
	AccessToken         string `json:"access_token"`
	BotUserID           string `json:"bot_user_id"`
	Enterprise          string `json:"enterprise"`
	IsEnterpriseInstall bool   `json:"is_enterprise_install"`
}

type SlackOauthClient struct {
	*oauth2.Config
}

var slackConf *SlackOauthClient = &SlackOauthClient{
	Config: &oauth2.Config{
		ClientID:     os.Getenv("SLACK_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("SLACK_OAUTH_CLIENT_SECRET"),
		RedirectURL:  fmt.Sprintf("%s/oauth2/slack/callback", os.Getenv("APP_HOST")),
		Scopes: []string{
			"users.profile:read",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://slack.com/oauth/v2/authorize",
			TokenURL: "https://slack.com/api/oauth.v2.access",
		},
	},
}

func (s *SlackOauthClient) Exchange(code string) (*SlackOauthResponse, error) {
	form := url.Values{}
	form.Add("code", code)
	form.Add("client_id", slackConf.ClientID)
	form.Add("client_secret", slackConf.ClientSecret)
	resp, err := http.PostForm(slackConf.Endpoint.TokenURL, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("slack Oauth2 response status code is not 200: %d", resp.StatusCode)
	}
	var oauthResp SlackOauthResponse
	if err := json.NewDecoder(resp.Body).Decode(&oauthResp); err != nil {
		return nil, err
	}
	return &oauthResp, nil
}

type SlackUserProfile struct {
	Profile struct {
		RealName string `json:"real_name"`
		Email    string `json:"email"`
		Title    string `json:"title"`
	} `json:"profile"`
}

func (s *SlackOauthClient) GetUsersProfile(slackUserID, accessToken string) (*SlackUserProfile, error) {
	query := url.Values{}
	query.Add("user", slackUserID)
	req, err := http.NewRequest(http.MethodGet, "https://slack.com/api/users.profile.get", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.URL.RawQuery = query.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var up SlackUserProfile
	if err := json.NewDecoder(resp.Body).Decode(&up); err != nil {
		return nil, err
	}
	return &up, nil
}

// SlackOauth2CallbackHandler handles the callback from Slack OAuth2
// anyone who signs up with Slack will be redirected to this handler
// and will be automatically created with the "Recruiter" role
func SlackOauth2CallbackHandler(c echo.Context) error {
	if slackOauth2StateToken != c.QueryParam("state") {
		// the state token is invalid, someone may be trying to intercept our login flow
		return echo.ErrUnauthorized
	}
	if c.QueryParams().Has("error") {
		return c.String(http.StatusOK, "You denied access to your Slack account")
	}
	// Exchange code for token
	code := c.QueryParam("code")
	oauthResp, err := slackConf.Exchange(code)
	if err != nil {
		return err
	}
	// Get user profile
	slackUserProfile, err := slackConf.GetUsersProfile(oauthResp.AuthedUser.ID, oauthResp.AccessToken)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"profile": slackUserProfile,
		"oauth":   oauthResp,
	})
}

func RecruiterLoginHanlder(c echo.Context) error {
	slackOauth2StateToken = randutil.GenerateOauthStateToken()
	url := slackConf.AuthCodeURL(slackOauth2StateToken)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}
