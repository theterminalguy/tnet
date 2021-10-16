package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/workexperience"
	"github.com/10hourlabs/tentn/util/date"
	"github.com/google/uuid"
)

type WorkExperienceRepository struct{}

type WorkExperienceParams struct {
	ApplicantUUID uuid.UUID `json:"applicant_uuid" validate:"required"`
	CompanyName   string    `json:"company_name" validate:"required"`
	Location      string    `json:"location" validate:"required"`
	JobTitle      string    `json:"job_title" validate:"required"`
	Description   string    `json:"description" validate:"required"`
	StartDate     string    `json:"start_date" validate:"datetime=2006-01-02T15:04:05Z07:00"`
	EndDate       time.Time `json:"end_date"`
}

func NewWorkExperienceRepository() *WorkExperienceRepository {
	return &WorkExperienceRepository{}
}

func (*WorkExperienceRepository) GetAll() ([]*ent.WorkExperience, error) {
	records, err := dBConn.WorkExperience.Query().
		Where(workexperience.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*WorkExperienceRepository) GetByUUID(id uuid.UUID) (*ent.WorkExperience, error) {
	record, err := dBConn.WorkExperience.Query().
		Where(workexperience.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, RecordNotFoundError
	}
	return record, nil
}

func (*WorkExperienceRepository) Create(p WorkExperienceParams) (*ent.WorkExperience, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	a, err := NewApplicantRepository().GetByUUID(p.ApplicantUUID)
	if err != nil {
		return nil, err
	}
	sd, err := date.ToRFC3339(p.StartDate)
	if err != nil {
		return nil, err
	}
	// TODO: fix bug with end date
	//ed, _ := date.ToRFC3339(p.EndDate)
	record, err := dBConn.WorkExperience.
		Create().
		SetApplicantID(a.ID).
		SetCompanyName(p.CompanyName).
		SetLocation(p.Location).
		SetJobTitle(p.JobTitle).
		SetDescription(p.Description).
		SetStartDate(*sd).
		SetNillableEndDate(nil).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *WorkExperienceRepository) Update(id uuid.UUID, p WorkExperienceParams) (*ent.WorkExperience, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := r.GetByUUID(id)
	if err != nil {
		return nil, err
	}
	// TODO, fix date
	sd, err := date.ToRFC3339(p.StartDate)
	if err != nil {
		return nil, err
	}
	_, err = record.Update().
		SetCompanyName(p.CompanyName).
		SetLocation(p.Location).
		SetJobTitle(p.JobTitle).
		SetDescription(p.Description).
		SetStartDate(*sd).
		SetNillableEndDate(nil).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (r *WorkExperienceRepository) DeleteByUUID(id uuid.UUID) error {
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
