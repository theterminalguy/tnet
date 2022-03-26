package repository

import (
	"errors"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/jobapplication"
	"github.com/10hourlabs/tentn/ent/predicate"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/google/uuid"
)

type JobApplicationQuerier interface {
	GetAllForTalent(talentID int) ([]*ent.JobApplication, error)
	GetAll() ([]*ent.JobApplication, error)
	GetByID(id uuid.UUID) (*ent.JobApplication, error)
	Create(p JobApplicationParams) (*ent.JobApplication, error)
	Update(id uuid.UUID, p JobApplicationParams) (*ent.JobApplication, []error)
	DeleteByID(id uuid.UUID) error
}

type JobApplicationRepository struct{}

type JobApplicationParams struct {
	JobUUID        uuid.UUID `json:"job_uuid" validate:"required"`
	TalentID       uuid.UUID `json:"talent_id" validate:"required"`
	ReferralSource string    `json:"referral_source"`

	// The Note field is only used internal by 10HL admins/recruiters
	// this is used to keep track of useful notes related to a candidate
	// Talent for a specific job
	Note   string `json:"note"`
	Status string `json:"status"`
}

func NewJobApplicationRepository() *JobApplicationRepository {
	return &JobApplicationRepository{}
}

func (*JobApplicationRepository) Filter(prd ...predicate.JobApplication) ([]*ent.JobApplication, error) {
	jobApplications, err := dBConn.JobApplication.Query().
		Where(prd...).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return jobApplications, nil
}

func (*JobApplicationRepository) GetAllForTalent(talentID uuid.UUID) ([]*ent.JobApplication, error) {
	records, err := dBConn.JobApplication.Query().
		Where(jobapplication.And(
			jobapplication.TalentID(talentID),
			jobapplication.DeletedAtIsNil())).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
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

func (*JobApplicationRepository) GetByID(id uuid.UUID) (*ent.JobApplication, error) {
	record, err := dBConn.JobApplication.Query().
		Where(jobapplication.ID(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*JobApplicationRepository) Create(p JobApplicationParams) (*ent.JobApplication, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	a, err := NewTalentRepository().GetByID(p.TalentID)
	if err != nil {
		return nil, err
	}
	j, err := NewJobRepository().GetByID(p.JobUUID)
	if err != nil {
		return nil, err
	}
	record, err := dBConn.JobApplication.
		Create().
		SetTalentID(a.ID).
		SetJobID(j.ID).
		SetReferralSource(p.ReferralSource).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *JobApplicationRepository) Update(id uuid.UUID, p JobApplicationParams) (*ent.JobApplication, []error) {
	err := validateParams(p, "TalentID")
	if err != nil {
		return nil, []error{err}
	}

	err = validateParams(p, "JobUUID")
	if err != nil {
		return nil, []error{err}
	}

	record, err := r.GetByID(id)
	if err != nil {
		return nil, []error{err}
	}

	var vldErrs []error
	bldr := record.Update()

	// Set and Validate Note if provided
	if vldErr := setNillableStringField(p.Note, func(v string) error {
		err := validateParams(p, "Note")
		if err != nil {
			return err
		}
		bldr.SetNote(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Status if provided
	if vldErr := setNillableStringField(p.Status, func(v string) error {
		err := validateParams(p, "Status")
		if err != nil {
			return err
		}
		bldr.SetStatus(jobapplication.Status(v))
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	if collection.HasAny(vldErrs) {
		return nil, vldErrs
	}

	record, err = bldr.Save(dBContext)
	if err != nil {
		return nil, []error{err}
	}
	return record, nil
}

func (r *JobApplicationRepository) DeleteByID(id uuid.UUID) error {
	record, err := r.GetByID(id)
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

func (*JobApplicationRepository) AlreadyApplied(jobID, talentID uuid.UUID) error {
	records, err := dBConn.JobApplication.Query().Where(
		jobapplication.And(
			jobapplication.JobID(jobID),
			jobapplication.TalentID(talentID),
		)).All(dBContext)
	if err != nil {
		return err
	}
	if collection.HasAny(records) {
		return errors.New("talent already applied for this job")
	}

	return nil
}
