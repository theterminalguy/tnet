package handler

import (
	"fmt"
	"os"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
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

func NewJWTClaims(u *ent.User, c echo.Context) *JWTClaims {
	return &JWTClaims{
		Role:      string(u.Role),
		TokenType: IDToken,
		Approved:  u.Approved,
		StandardClaims: jwt.StandardClaims{
			Audience: c.Request().Host,
			// TODO: token expiration should be configurable
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    c.Request().Host,
			Subject:   u.UUID.String(),
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
