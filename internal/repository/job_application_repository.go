package repository

import (
	"errors"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/jobapplication"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/google/uuid"
)

type JobApplicationRepository struct{}

type JobApplicationParams struct {
	JobUUID        uuid.UUID `json:"job_uuid" validate:"required"`
	ApplicantUUID  uuid.UUID `json:"applicant_uuid" validate:"required"`
	ReferralSource string    `json:"referral_source"`

	// The Note field is only used internal by 10HL admins/recruiters
	// this is used to keep track of useful notes related to a candidate
	// application for a specific job
	Note   string `json:"note"`
	Status string `json:"status"`
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
	a, err := NewApplicantRepository().GetByUUID(p.ApplicantUUID)
	if err != nil {
		return nil, err
	}
	j, err := NewJobRepository().GetByUUID(p.JobUUID)
	if err != nil {
		return nil, err
	}
	records, err := dBConn.JobApplication.Query().Where(
		jobapplication.And(
			jobapplication.JobID(j.ID),
			jobapplication.ApplicantID(a.ID),
		)).All(dBContext)
	if collection.HasAny(records) {
		return nil, errors.New("applicant already applied for this job")
	}
	record, err := dBConn.JobApplication.
		Create().
		SetApplicantID(a.ID).
		SetJobID(j.ID).
		SetReferralSource(p.ReferralSource).
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
	_, err = record.Update().
		SetNote(p.Note).
		SetStatus(jobapplication.Status(p.Status)).
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
	_, err = record.Update().
		SetDeletedAt(time.Now()).
		Save(dBContext)
	if err != nil {
		return err
	}
	return nil
}
