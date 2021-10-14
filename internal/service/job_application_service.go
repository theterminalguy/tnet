package service

import (
	"github.com/google/uuid"
)

type JobApplicationService struct{}

func NewJobApplicationService() *JobApplicationService {
	return &JobApplicationService{}
}

func (*JobApplicationService) Apply(jobUUID, applicantUUID uuid.UUID) {
	// TODO
	// user must have a linkedin profile
	// user must have a GitHub profile for Enginering role
}
