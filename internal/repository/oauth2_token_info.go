package repository

import (
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/google/uuid"
	"github.com/theterminalguy/tenlog"
	"github.com/theterminalguy/tnet/ent"
)

type Oauth2TokenInfo struct {
	*ent.Oauth2Token
}

func (o Oauth2TokenInfo) New() oauth2.TokenInfo {
	return Oauth2TokenInfo{}
}

func (o Oauth2TokenInfo) GetClientID() string {
	return o.Oauth2Token.Oauth2ClientID.String()
}

func (o Oauth2TokenInfo) SetClientID(clientID string) {
	o.Oauth2Token.Oauth2ClientID = uuid.MustParse(clientID)
}

func (o Oauth2TokenInfo) GetUserID() string {
	tenlog.Debug("Getting user id", o.Oauth2Token.UserID)
	return o.Oauth2Token.UserID.String()
}

func (o Oauth2TokenInfo) SetUserID(userID string) {
	o.Oauth2Token.UserID = uuid.MustParse(userID)
}

func (o Oauth2TokenInfo) GetRedirectURI() string {
	return o.Oauth2Token.RedirectURI
}

func (o Oauth2TokenInfo) SetRedirectURI(redirectURI string) {
	o.Oauth2Token.RedirectURI = redirectURI
}

func (o Oauth2TokenInfo) GetScope() string {
	return o.Oauth2Token.Scopes
}

func (o Oauth2TokenInfo) SetScope(scope string) {
	o.Oauth2Token.Scopes = scope
}

func (o Oauth2TokenInfo) GetCode() string {
	return o.Oauth2Token.Code
}

func (o Oauth2TokenInfo) SetCode(code string) {
	o.Oauth2Token.Code = code
}

func (o Oauth2TokenInfo) GetCodeCreateAt() time.Time {
	return o.Oauth2Token.CodeCreatedAt
}

func (o Oauth2TokenInfo) SetCodeCreateAt(codeCreateAt time.Time) {
	o.Oauth2Token.CodeCreatedAt = codeCreateAt
}

func (o Oauth2TokenInfo) GetCodeExpiresIn() time.Duration {
	return time.Duration(o.Oauth2Token.CodeExpiresIn)
}

func (o Oauth2TokenInfo) SetCodeExpiresIn(codeExpiresIn time.Duration) {
	o.Oauth2Token.CodeExpiresIn = int64(codeExpiresIn)
}

func (o Oauth2TokenInfo) GetCodeChallenge() string {
	return o.Oauth2Token.CodeChallenge
}

func (o Oauth2TokenInfo) SetCodeChallenge(codeChallenge string) {
	o.Oauth2Token.CodeChallenge = codeChallenge
}

func (o Oauth2TokenInfo) GetCodeChallengeMethod() oauth2.CodeChallengeMethod {
	return oauth2.CodeChallengeMethod(o.Oauth2Token.CodeChallengeMethod)
}

func (o Oauth2TokenInfo) SetCodeChallengeMethod(codeChallengeMethod oauth2.CodeChallengeMethod) {
	o.Oauth2Token.CodeChallengeMethod = string(codeChallengeMethod)
}

func (o Oauth2TokenInfo) GetAccess() string {
	return o.Oauth2Token.AccessToken
}

func (o Oauth2TokenInfo) SetAccess(access string) {
	o.Oauth2Token.AccessToken = access
}

func (o Oauth2TokenInfo) GetAccessCreateAt() time.Time {
	return o.Oauth2Token.AccessTokenCreatedAt
}

func (o Oauth2TokenInfo) SetAccessCreateAt(accessCreateAt time.Time) {
	o.Oauth2Token.AccessTokenCreatedAt = accessCreateAt
}

func (o Oauth2TokenInfo) GetAccessExpiresIn() time.Duration {
	return time.Duration(o.Oauth2Token.AccessTokenExpiresIn)
}

func (o Oauth2TokenInfo) SetAccessExpiresIn(accessExpiresIn time.Duration) {
	o.Oauth2Token.AccessTokenExpiresIn = int64(accessExpiresIn)
}

func (o Oauth2TokenInfo) GetRefresh() string {
	return o.Oauth2Token.RefreshToken
}

func (o Oauth2TokenInfo) SetRefresh(refresh string) {
	o.Oauth2Token.RefreshToken = refresh
}

func (o Oauth2TokenInfo) GetRefreshCreateAt() time.Time {
	return o.Oauth2Token.RefreshTokenCreatedAt
}

func (o Oauth2TokenInfo) SetRefreshCreateAt(refreshCreateAt time.Time) {
	o.Oauth2Token.RefreshTokenCreatedAt = refreshCreateAt
}

func (o Oauth2TokenInfo) GetRefreshExpiresIn() time.Duration {
	return time.Duration(o.Oauth2Token.RefreshTokenExpiresIn)
}

func (o Oauth2TokenInfo) SetRefreshExpiresIn(refreshExpiresIn time.Duration) {
	o.Oauth2Token.RefreshTokenExpiresIn = int64(refreshExpiresIn)
}
