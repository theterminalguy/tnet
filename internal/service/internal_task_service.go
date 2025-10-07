package service

import (
	"fmt"

	"github.com/theterminalguy/tentn/ent"
	repo "github.com/theterminalguy/tentn/internal/repository"
	"github.com/theterminalguy/tentn/util"
)

type TaskRunner func(name, params, executor, password string) error

type InternalTaskService struct {
	TaskRepo *repo.InternalTaskRepository
}

func NewInternalTaskService() *InternalTaskService {
	return &InternalTaskService{}
}

func (it *InternalTaskService) Create(runTask TaskRunner, t repo.InternalTaskParams) (*ent.InternalTask, error) {
	err := runTask(t.Name, util.MapToStringParams(t.Params), t.Executor, t.Password)
	if err != nil {
		tr, tErr := it.TaskRepo.Create(repo.InternalTaskParams{
			Name:      t.Name,
			Params:    t.Params,
			Executor:  t.Executor,
			Succeeded: false,
			Error:     err.Error(),
		})
		return tr, fmt.Errorf("%s: %s", tErr, err)
	}
	tr, tErr := it.TaskRepo.Create(repo.InternalTaskParams{
		Name:      t.Name,
		Params:    t.Params,
		Executor:  t.Executor,
		Succeeded: true,
		Error:     "",
	})
	if tErr != nil {
		return nil, tErr
	}
	return tr, nil
}
