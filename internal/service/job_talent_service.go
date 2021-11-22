package service

import (
	"github.com/google/uuid"
)

type JobTalentService struct{}

func NewJobTalentService() *JobTalentService {
	return &JobTalentService{}
}

func (*JobTalentService) Apply(jobUUID, applicantUUID uuid.UUID) {
	// TODO
	// user must have a linkedin profile
	// user must have a GitHub profile for Enginering role
}
