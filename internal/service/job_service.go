package service

import (
	"github.com/10hourlabs/tentn/internal/repo"
)

type JobService struct {
	JobRepo *repo.JobRepository
}

func NewJobService() *JobService {
	return &JobService{
		JobRepo: repo.NewJobRepository(),
	}
}
