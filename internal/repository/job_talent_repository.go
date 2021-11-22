package repository

import (
	"errors"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/jobtalent"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/google/uuid"
)

type JobTalentRepository struct{}

type JobTalentParams struct {
	JobUUID        uuid.UUID `json:"job_uuid" validate:"required"`
	TalentUUID     uuid.UUID `json:"talent_uuid" validate:"required"`
	ReferralSource string    `json:"referral_source"`

	// The Note field is only used internal by 10HL admins/recruiters
	// this is used to keep track of useful notes related to a candidate
	// Talent for a specific job
	Note   string `json:"note"`
	Status string `json:"status"`
}

func NewJobTalentRepository() *JobTalentRepository {
	return &JobTalentRepository{}
}

func (*JobTalentRepository) GetAll() ([]*ent.JobTalent, error) {
	records, err := dBConn.JobTalent.Query().
		Where(jobtalent.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*JobTalentRepository) GetByUUID(id uuid.UUID) (*ent.JobTalent, error) {
	record, err := dBConn.JobTalent.Query().
		Where(jobtalent.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, RecordNotFoundError
	}
	return record, nil
}

func (*JobTalentRepository) Create(p JobTalentParams) (*ent.JobTalent, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	a, err := NewTalentRepository().GetByUUID(p.TalentUUID)
	if err != nil {
		return nil, err
	}
	j, err := NewJobRepository().GetByUUID(p.JobUUID)
	if err != nil {
		return nil, err
	}
	records, err := dBConn.JobTalent.Query().Where(
		jobtalent.And(
			jobtalent.JobID(j.ID),
			jobtalent.TalentID(a.ID),
		)).All(dBContext)
	if collection.HasAny(records) {
		return nil, errors.New("Talent already applied for this job")
	}
	record, err := dBConn.JobTalent.
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

func (r *JobTalentRepository) Update(id uuid.UUID, p JobTalentParams) (*ent.JobTalent, []error) {
	err := validateParams(p, "TalentUUID")
	if err != nil {
		return nil, []error{err}
	}
	record, err := r.GetByUUID(id)
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
		bldr.SetStatus(jobtalent.Status(v))
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

func (r *JobTalentRepository) DeleteByUUID(id uuid.UUID) error {
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
