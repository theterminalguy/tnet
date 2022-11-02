package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/10hourlabs/tentn/internal/entsm"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/internal/tokgen"
	"github.com/10hourlabs/tentn/randutil"
	"github.com/10hourlabs/tentn/util/photo"
	"github.com/gorilla/sessions"
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
	Error               string `json:"error"`
	AppID               string `json:"app_id"`
	Scope               string `json:"scope"`
	TokenType           string `json:"token_type"`
	AccessToken         string `json:"access_token"`
	BotUserID           string `json:"bot_user_id"`
	Enterprise          string `json:"enterprise"`
	IsEnterpriseInstall bool   `json:"is_enterprise_install"`
}

func (s *SlackOauthResponse) GetRecruitersEmail() string {
	return fmt.Sprintf("%s@%s.slack.com", s.AuthedUser.ID, s.Team.ID)
}

type SlackOauth2Client struct {
	*oauth2.Config
}

var slackConf *SlackOauth2Client = &SlackOauth2Client{
	Config: &oauth2.Config{
		ClientID:     os.Getenv("SLACK_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("SLACK_OAUTH_CLIENT_SECRET"),
		RedirectURL:  fmt.Sprintf("%s/oauth2/slack/callback", os.Getenv("APP_HOST")),
		Scopes: []string{
			"users.profile:read",
			"channels:history",
			"channels:join",
			"channels:manage",
			"chat:write",
			"conversations.connect:write",
			"files:write",
			"groups:history",
			"groups:write",
			"im:history",
			"mpim:history",
			"team.billing:read",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://slack.com/oauth/v2/authorize",
			TokenURL: "https://slack.com/api/oauth.v2.access",
		},
	},
}

func (s *SlackOauth2Client) Exchange(code string) (*SlackOauthResponse, error) {
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
	if !oauthResp.Ok {
		return nil, fmt.Errorf("slack Oauth2 response is not ok: %v", oauthResp.Error)
	}
	return &oauthResp, nil
}

type SlackUserProfile struct {
	Profile struct {
		FirstName string  `json:"first_name"`
		LastName  *string `json:"last_name"`
		PhotoURL  string  `json:"image_192"`
		Email     string  `json:"email"`
		Title     string  `json:"title"`
		Phone     string  `json:"phone"`
	} `json:"profile"`
}

func (s *SlackUserProfile) GetLastName(fname string) string {
	if s.Profile.LastName == nil {
		return fname
	}
	return *s.Profile.LastName
}

func (s *SlackUserProfile) FullName() string {
	fname := s.Profile.FirstName
	return fmt.Sprintf("%s %s", fname, s.GetLastName(fname))
}

func (s *SlackUserProfile) GetPhotoURL() string {
	if s.Profile.PhotoURL == "" || strings.Contains(s.Profile.PhotoURL, "gravatar") {
		fname := s.Profile.FirstName
		return photo.GenerateDefaultPhoto(fname, s.GetLastName(fname))
	}
	return s.Profile.PhotoURL
}

func (s *SlackOauth2Client) GetUsersProfile(slackUserID, accessToken string) (*SlackUserProfile, error) {
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

var (
	session *sessions.Session
)

// SlackOauth2CallbackHandler handles the callback from Slack OAuth2
// anyone who signs up with Slack will be redirected to this handler
// and will be automatically created with the "Recruiter" role
func SlackOauth2CallbackHandler(c echo.Context) error {
	if slackOauth2StateToken != c.QueryParam("state") {
		// the state token is invalid, someone may be trying to intercept our login flow
		return echo.ErrUnauthorized
	}
	if err := c.QueryParam("error"); err != "" {
		// TODO: return a better error message
		return c.String(http.StatusOK, fmt.Sprintf("You denied access to your Slack account: %s", err))
	}
	// Exchange code for token
	code := c.QueryParam("code")
	oauthResp, err := slackConf.Exchange(code)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Failed to exchange code for token: %s", err))
	}
	// Get user slack profile
	slackUserProfile, err := slackConf.GetUsersProfile(oauthResp.AuthedUser.ID, oauthResp.AccessToken)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Failed to get user profile: %s", err))
	}
	rs := service.NewRecruiterService()
	u, err := rs.InstallSlackApp(
		repo.UserParams{
			FirstName: slackUserProfile.Profile.FirstName,
			LastName:  slackUserProfile.GetLastName(slackUserProfile.Profile.FirstName),
			PhotoURL:  slackUserProfile.GetPhotoURL(),
			Email:     oauthResp.GetRecruitersEmail(),
		},
		repo.SlackAppInstallParams{
			TeamID:              oauthResp.Team.ID,
			TeamName:            oauthResp.Team.Name,
			AuthedUserID:        oauthResp.AuthedUser.ID,
			AuthedUserEmail:     slackUserProfile.Profile.Email,
			AuthedUserTitle:     slackUserProfile.Profile.Title,
			AuthedUserPhone:     slackUserProfile.Profile.Phone,
			AppID:               oauthResp.AppID,
			BotUserID:           oauthResp.BotUserID,
			AccessToken:         oauthResp.AccessToken,
			TokenType:           oauthResp.TokenType,
			Scope:               oauthResp.Scope,
			IsEnterpriseInstall: oauthResp.IsEnterpriseInstall,
		})
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Failed to install Slack app: %s", err))
	}
	ses, err := repo.NewSessionRepository().GetSessionByTeamID(oauthResp.Team.ID)
	if ses != nil && ses.DeletedAt == nil {
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("https://slack.com/app_redirect?app=%s", oauthResp.AppID))
	}
	if err != nil {
		store := entsm.GetSessionStore()
		session, _ = store.Get(c.Request(), entsm.DefaultSessionName)
		session.Values["authenticated"] = true
		session.Values["team_id"] = oauthResp.Team.ID
	}
	claims := tokgen.NewJWTClaims(u, &tokgen.JWTMeta{
		Audience: c.Request().Host,
		Issuer:   c.Request().Host,
	})
	token, err := claims.GenerateToken()
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Failed to generate token: %s", err))
	}
	session.Values["token"] = token
	session.Options.Domain = os.Getenv("APP_HOST")
	session.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("https://slack.com/app_redirect?app=%s", oauthResp.AppID))
}

func RecruiterLoginHanlder(c echo.Context) error {
	store := entsm.GetSessionStore()
	session, _ := store.Get(c.Request(), entsm.DefaultSessionName)
	if _, ok := session.Values["authenticated"].(bool); ok {
		return c.Redirect(http.StatusTemporaryRedirect, os.Getenv("APP_HOST"))
	}
	slackOauth2StateToken = randutil.GenerateOauthStateToken()
	url := slackConf.AuthCodeURL(slackOauth2StateToken)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func LogoutHandler(c echo.Context) error {
	store := entsm.GetSessionStore()
	session, _ := store.Get(c.Request(), entsm.DefaultSessionName)
	// Revoke users authentication
	session.Options.MaxAge = -1
	session.Save(c.Request(), c.Response())
	return c.JSON(http.StatusOK, "logged out")
}
