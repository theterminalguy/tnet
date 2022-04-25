package task

import (
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/google/uuid"
)

type TaskApproveClient struct {
	Oauth2Repo    repo.Oauth2ClientRepository
	Oauth2Service service.Oauth2ClientService
}

func NewTaskApproveClient() *TaskApproveClient {
	return &TaskApproveClient{
		Oauth2Repo:    *repo.NewOauth2ClientRepository(),
		Oauth2Service: *service.NewOauth2ClientService(),
	}
}

func (t *TaskApproveClient) Run(clientID string) error {
	client, err := t.Oauth2Repo.GetByUUID(uuid.MustParse(clientID))
	if err != nil {
		return err
	}
	err = t.Oauth2Service.ApproveClient(*client)
	if err != nil {
		return err
	}
	return nil
}
