package repository

import (
	"github.com/10hourlabs/tentn/ent"
	"github.com/google/uuid"
)

type JobFileUploadRepository struct{}

type JobFileUploadParams struct {
	FileUrl string `json:"file_url" validate:"required"`
}

func NewJobFileUploadRepository() *JobFileUploadRepository {
	return &JobFileUploadRepository{}
}

func (*JobFileUploadRepository) GetAll() ([]*ent.JobFileUpload, error) {
	return nil, nil
}

func (*JobFileUploadRepository) GetByID() (*ent.JobFileUpload, error) {
	return nil, nil
}

func (*JobFileUploadRepository) Create(request JobFileUploadParams) (*ent.JobFileUpload, error) {
	err := ValidateParams(request)
	if err != nil {
		return nil, err
	}

	result, err := dBConn.JobFileUpload.
		Create().
		SetFileURL(request.FileUrl).
		Save(dBContext)

	if err != nil {
		return nil, err
	}

	return result, err
}

func (*JobFileUploadRepository) Update(id uuid.UUID, param *JobFileUploadParams) (*ent.JobFileUpload, error) {
	return nil, nil
}

func (*JobFileUploadRepository) DeleteByID(id uuid.UUID) error {
	return nil
}
