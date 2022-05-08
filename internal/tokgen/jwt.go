package tokgen

import (
	"fmt"
	"os"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/golang-jwt/jwt"
)

type Token string

const (
	IDToken     Token = "id_token"
	AccessToken Token = "access_token"
)

// SignedKeyID is the key id for the signing key
// See https://stackoverflow.com/questions/43867440/whats-the-meaning-of-the-kid-claim-in-a-jwt-token
const DefaultSigningKeyID = ""
const DefaultAccessTokenExp = time.Hour * 24 * 7   // 7 days
const DefaultRefreshTokenExp = time.Hour * 24 * 30 // 30 days

func GetDefualtSigningKey() []byte {
	return []byte(os.Getenv("JWT_SIGNING_KEY"))
}

func GetDefaultSigningMethod() jwt.SigningMethod {
	return jwt.SigningMethodHS256
}

type JWTClaims struct {
	TokenType Token `json:"token_type"`
	jwt.StandardClaims
}

type JWTMeta struct {
	Audience string
	Issuer   string
}

func NewJWTClaims(u *ent.User, m *JWTMeta) *JWTClaims {
	return &JWTClaims{
		TokenType: IDToken,
		StandardClaims: jwt.StandardClaims{
			Audience: m.Audience,
			// TODO: token expiration should be configurable, currently set to 1 day
			ExpiresAt: time.Now().Add(DefaultAccessTokenExp).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    m.Issuer,
			Subject:   u.ID.String(),
		},
	}
}

func (c *JWTClaims) GenerateToken() (string, error) {
	jwttok := jwt.NewWithClaims(GetDefaultSigningMethod(), c)
	token, err := jwttok.SignedString(GetDefualtSigningKey())
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}
	return token, nil
}
