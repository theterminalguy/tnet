package repository

import (
	"github.com/10hourlabs/tentn/ent"
	"github.com/google/uuid"
)

type FileUploadRepository struct{}

type FileUploadParams struct {
	FileUrl string `json:"file_url" validate:"required"`
}

func NewFileUploadRepository() *FileUploadRepository {
	return &FileUploadRepository{}
}

func (*FileUploadRepository) GetAll() ([]*ent.FileUpload, error) {
	return nil, nil
}

func (*FileUploadRepository) GetByID() (*ent.FileUpload, error) {
	return nil, nil
}

func (*FileUploadRepository) Create(request FileUploadParams) (*ent.FileUpload, error) {
	err := ValidateParams(request)
	if err != nil {
		return nil, err
	}

	result, err := dBConn.FileUpload.
		Create().
		SetFileURL(request.FileUrl).
		Save(dBContext)

	if err != nil {
		return nil, err
	}

	return result, err
}

func (*FileUploadRepository) Update(id uuid.UUID, param *FileUploadParams) (*ent.FileUpload, error) {
	return nil, nil
}

func (*FileUploadRepository) DeleteByID(id uuid.UUID) error {
	return nil
}
