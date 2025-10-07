package repository

import (
	"strings"

	"github.com/theterminalguy/tentn/ent"
	"golang.org/x/crypto/bcrypt"
)

type Oauth2ClientInfo struct {
	*ent.Oauth2Client
}

func (c *Oauth2ClientInfo) GetID() string {
	return c.ID.String()
}

func (c *Oauth2ClientInfo) GetSecret() string {
	return c.HashedSecret
}

func (c *Oauth2ClientInfo) GetDomain() string {
	return strings.Join(c.RedirectUris, ",")
}

func (c *Oauth2ClientInfo) GetUserID() string {
	return c.UserID.String()
}

func (c *Oauth2ClientInfo) VerifyPassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(c.HashedSecret), []byte(password)); err == nil {
		return true
	}
	return false
}
