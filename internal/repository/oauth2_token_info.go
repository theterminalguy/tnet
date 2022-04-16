package repository

import (
	"github.com/10hourlabs/tentn/ent"
	"github.com/google/uuid"
)

/*TokenInfo interface {
	GetUserID() string
	SetUserID(string)
	GetRedirectURI() string
	SetRedirectURI(string)
	GetScope() string
	SetScope(string)

	GetCode() string
	SetCode(string)
	GetCodeCreateAt() time.Time
	SetCodeCreateAt(time.Time)
	GetCodeExpiresIn() time.Duration
	SetCodeExpiresIn(time.Duration)
	GetCodeChallenge() string
	SetCodeChallenge(string)
	GetCodeChallengeMethod() CodeChallengeMethod
	SetCodeChallengeMethod(CodeChallengeMethod)

	GetAccess() string
	SetAccess(string)
	GetAccessCreateAt() time.Time
	SetAccessCreateAt(time.Time)
	GetAccessExpiresIn() time.Duration
	SetAccessExpiresIn(time.Duration)

	GetRefresh() string
	SetRefresh(string)
	GetRefreshCreateAt() time.Time
	SetRefreshCreateAt(time.Time)
	GetRefreshExpiresIn() time.Duration
	SetRefreshExpiresIn(time.Duration)
}*/

type Oauth2TokenInfo struct {
	*ent.Oauth2Token
}

func (o *Oauth2TokenInfo) New() *Oauth2TokenInfo {
	return &Oauth2TokenInfo{}
}

func (o *Oauth2TokenInfo) GetClientID() string {
	return o.Oauth2Token.Oauth2ClientID.String()
}

func (o *Oauth2TokenInfo) SetClientID(clientID string) {
	o.Oauth2Token.Oauth2ClientID = uuid.MustParse(clientID)
}
