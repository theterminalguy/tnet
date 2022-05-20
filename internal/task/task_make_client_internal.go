package task

import (
	"fmt"

	"github.com/10hourlabs/tentn/ent/oauth2client"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/util"
	"github.com/google/uuid"
)

type TaskMakeClientInternal struct {
	Oauth2Repo repo.Oauth2ClientRepository
}

func NewTaskMakeClientInternal() *TaskMakeClientInternal {
	return &TaskMakeClientInternal{
		Oauth2Repo: *repo.NewOauth2ClientRepository(),
	}
}

func (t *TaskMakeClientInternal) Run(params string) error {
	m := util.StringParamsToMap(params)
	clientID, ok := m["client_id"]
	if !ok {
		return fmt.Errorf("provide a client_id")
	}
	client, err := t.Oauth2Repo.GetByUUID(uuid.MustParse(clientID))
	if err != nil {
		return err
	}
	if client.IsInternal {
		return nil
	}
	err = t.Oauth2Repo.UpdateFields(client, map[string]interface{}{
		oauth2client.FieldIsInternal: true,
	})
	if err != nil {
		return err
	}
	return nil
}
