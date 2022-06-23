package repository

import (
	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/jobcollection"
	"github.com/google/uuid"
)

type JobCollectionRepository struct{}

type JobCollectionParams struct {
	RecruiterID uuid.UUID
	Title       string `json:"title" validate:"required"`
	Status      string `json:"status" validate:"required"`
}

func NewJobCollectionRepository() *JobCollectionRepository {
	return &JobCollectionRepository{}
}

func (*JobCollectionRepository) GetAll(recruiterID uuid.UUID) ([]*ent.JobCollection, error) {
	return nil, nil
}

func (*JobCollectionRepository) GetByRecruiterID(recruiterID uuid.UUID) ([]*ent.JobCollection, error) {
	j, err := dBConn.JobCollection.Query().
		Where(jobcollection.RecruiterID(recruiterID)).
		Where(jobcollection.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}

	return j, nil
}

func (*JobCollectionRepository) Create(request JobCollectionParams) (*ent.JobCollection, error) {
	err := ValidateParams(request)
	if err != nil {
		return nil, err
	}

	result, err := dBConn.JobCollection.
		Create().
		SetRecruiterID(request.RecruiterID).
		SetTitle(request.Title).
		SetStatus(jobcollection.Status(request.Status)).
		Save(dBContext)

	if err != nil {
		return nil, err
	}

	return result, err
}

func (*JobCollectionRepository) Update(id uuid.UUID, request JobCollectionParams) (*ent.JobCollection, error) {
	return nil, nil
}

func (*JobCollectionRepository) GetByID(id uuid.UUID) (*ent.JobCollection, error) {
	return nil, nil
}

func (*JobCollectionRepository) DeleteByID(id uuid.UUID) error {
	return nil
}
