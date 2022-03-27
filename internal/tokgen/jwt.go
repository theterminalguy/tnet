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
	IDToken     = "id_token"
	AccessToken = "access_token"
)

type JWTClaims struct {
	Role      string `json:"role"`
	TokenType Token  `json:"token_type"`
	Approved  bool   `json:"approved"`
	jwt.StandardClaims
}

type JWTMeta struct {
	Audience string
	Issuer   string
}

func NewJWTClaims(u *ent.User, m *JWTMeta) *JWTClaims {
	return &JWTClaims{
		Role:      string(u.Role),
		TokenType: IDToken,
		Approved:  u.Approved,
		StandardClaims: jwt.StandardClaims{
			Audience: m.Audience,
			// TODO: token expiration should be configurable, currently set to 1 day
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    m.Issuer,
			Subject:   u.ID.String(),
		},
	}
}

func (c *JWTClaims) GenerateToken() (string, error) {
	jwttok := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token, err := jwttok.SignedString([]byte(os.Getenv("JWT_SIGNED_SECRET")))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}
	return token, nil
}
