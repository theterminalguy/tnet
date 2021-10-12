package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/jobapplication"
	"github.com/google/uuid"
)

type JobApplicationRepository struct{}

type JobApplicationParams struct {

}

func NewJobApplicationRepository() *JobApplicationRepository {
	return &JobApplicationRepository{}
}

func (*JobApplicationRepository) GetAll() ([]*ent.JobApplication, error) {
	records, err := dBConn.JobApplication.Query().
		Where(jobapplication.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*JobApplicationRepository) GetByUUID(id uuid.UUID) (*ent.JobApplication, error) {
	record, err := dBConn.JobApplication.Query().
		Where(jobapplication.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, RecordNotFoundError
	}
	return record, nil
}

func (*JobApplicationRepository) Create(p JobApplicationParams) (*ent.JobApplication, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := dBConn.JobApplication.
		Create().
		// TODO: set other fields here
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *JobApplicationRepository) Update(id uuid.UUID, p JobApplicationParams) (*ent.JobApplication, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := r.GetByUUID(id)
	if err != nil {
		return nil, err
	}
	_, err = dBConn.JobApplication.Update().
		// TODO: set other fields here
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (r *JobApplicationRepository) DeleteByUUID(id uuid.UUID) error {
	record, err := r.GetByUUID(id)
	if err != nil {
		return err
	}
	_, err = dBConn.JobApplication.UpdateOne(record).
		SetDeletedAt(time.Now()).
		Save(dBContext)
	if err != nil {
		return err
	}
	return nil
}
