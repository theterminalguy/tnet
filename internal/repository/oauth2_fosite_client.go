package repository

import (
	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/oauth2client"
	"github.com/ory/fosite"
)

// https://github.com/ory/fosite/blob/0a48821b156f4a5dffa0f7149d30d5cf02636f37/client.go#L27
type Oauth2FositeClient struct {
	GrantTypes    []string
	ResponseTypes []string
	*ent.Oauth2Client
}

func (c Oauth2FositeClient) GetID() string {
	return c.ID.String()
}

func (c Oauth2FositeClient) GetHashedSecret() []byte {
	return []byte(c.HashedSecret)
}

func (c Oauth2FositeClient) GetRedirectURIs() []string {
	return c.RedirectUris
}

func (c Oauth2FositeClient) GetGrantTypes() fosite.Arguments {
	return fosite.Arguments(c.GrantTypes)
}

func (c Oauth2FositeClient) GetResponseTypes() fosite.Arguments {
	return fosite.Arguments(c.ResponseTypes)
}

func (c Oauth2FositeClient) GetScopes() fosite.Arguments {
	return fosite.Arguments(c.Scopes)
}

func (c Oauth2FositeClient) IsPublic() bool {
	return c.ClientType == oauth2client.ClientTypePublic
}

func (c Oauth2FositeClient) GetAudience() fosite.Arguments {
	return fosite.Arguments([]string{c.UserID.String()})
}
