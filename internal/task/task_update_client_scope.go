package task

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/theterminalguy/tnet/ent/oauth2client"
	repo "github.com/theterminalguy/tnet/internal/repository"
	"github.com/theterminalguy/tnet/util"
	"github.com/theterminalguy/tnet/util/collection"
)

type TaskUpdateClientScope struct {
	cr *repo.Oauth2ClientRepository
	sr *repo.Oauth2ScopeRepository
}

func NewTaskUpdateClientScope() *TaskUpdateClientScope {
	return &TaskUpdateClientScope{
		cr: repo.NewOauth2ClientRepository(),
		sr: repo.NewOauth2ScopeRepository(),
	}

}

func (t *TaskUpdateClientScope) Run(params string) error {
	m := util.StringParamsToMap(params)
	clientID, ok := m["client-id"]
	if !ok {
		return fmt.Errorf("provide a client ID")
	}
	scope, ok := m["scope"]
	if !ok {
		return fmt.Errorf("provide a scope")
	}
	scopes := strings.Split(scope, ",")
	for _, s := range scopes {
		if !t.sr.IsValid(s) {
			return fmt.Errorf("invalid scope %s", s)
		}
	}
	client, err := t.cr.GetByUUID(uuid.MustParse(clientID))
	if err != nil {
		return err
	}
	newScopes := collection.UniqueCombineStr(client.Scopes, scopes)
	return t.cr.UpdateFields(client, map[string]interface{}{
		oauth2client.FieldScopes: newScopes,
	})
}
