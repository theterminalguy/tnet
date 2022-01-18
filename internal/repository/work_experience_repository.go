package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/predicate"
	"github.com/10hourlabs/tentn/ent/workexperience"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/10hourlabs/tentn/util/date"
	"github.com/google/uuid"
)

type WorkExperienceQuerier interface {
	GetAll() ([]*ent.WorkExperience, error)
	GetAllForTalent(talentID int) ([]*ent.WorkExperience, error)
	GetByUUID(id uuid.UUID) (*ent.WorkExperience, error)
	Create(p WorkExperienceParams) (*ent.WorkExperience, error)
	Update(id uuid.UUID, p WorkExperienceParams) (*ent.WorkExperience, []error)
	DeleteByUUID(id uuid.UUID) error
}

type WorkExperienceRepository struct{}

type WorkExperienceParams struct {
	TalentUUID          uuid.UUID `json:"talent_uuid" validate:"required"`
	CompanyName         string    `json:"company_name" validate:"required"`
	Location            string    `json:"location" validate:"required"`
	JobTitle            string    `json:"job_title" validate:"required"`
	PrimaryTechnologies []string  `json:"primary_technologies" validate:"required"`
	Description         string    `json:"description" validate:"required"`
	StartDate           string    `json:"start_date" validate:"datetime=2006-01-02T15:04:05Z07:00"`
	EndDate             string    `json:"end_date"`
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

func (*WorkExperienceRepository) Filter(prd ...predicate.WorkExperience) ([]*ent.WorkExperience, error) {
	wkExps, err := dBConn.WorkExperience.Query().
		Where(prd...).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return wkExps, nil
}

func (*WorkExperienceRepository) GetAllForTalent(talentID int) ([]*ent.WorkExperience, error) {
	records, err := dBConn.WorkExperience.Query().
		Where(workexperience.And(
			workexperience.TalentID(talentID),
			workexperience.DeletedAtIsNil())).
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
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*WorkExperienceRepository) Create(p WorkExperienceParams) (*ent.WorkExperience, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}

	var sd *time.Time
	var ed *time.Time

	a, err := NewTalentRepository().GetByUUID(p.TalentUUID)
	if err != nil {
		return nil, err
	}
	sd, err = date.JSStringToRFC3339(p.StartDate)
	if err != nil {
		return nil, err
	}

	if p.EndDate != "" {
		ed, err = date.JSStringToRFC3339(p.EndDate)
		if err != nil {
			return nil, err
		}

		err = IsEqual(*sd, *ed)
		if err != nil {
			return nil, err
		}
	}
	// TODO: fix bug with end date
	//ed, _ := date.JSStringToRFC3339(p.EndDate)
	record, err := dBConn.WorkExperience.
		Create().
		SetTalentID(a.ID).
		SetCompanyName(p.CompanyName).
		SetLocation(p.Location).
		SetJobTitle(p.JobTitle).
		SetPrimaryTechnologies(p.PrimaryTechnologies).
		SetDescription(p.Description).
		SetStartDate(*sd).
		SetNillableEndDate(ed).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *WorkExperienceRepository) Update(id uuid.UUID, p WorkExperienceParams) (*ent.WorkExperience, []error) {
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

	sd := date.ToRFC3339(record.StartDate)

	setNillableStringField(p.StartDate, func(v string) error {
		sd = v
		return nil
	})

	period := &DatePeriod{
		StartDate: sd,
		EndDate:   p.EndDate,
	}

	if err = period.IsValid(func(startdate, enddate *time.Time) {
		bldr.SetStartDate(*startdate)
		bldr.SetNillableEndDate(enddate)
	}); err != nil {
		vldErrs = append(vldErrs, err)
	}

	// Set and Validate Location if provided
	if vldErr := setNillableStringField(p.Location, func(v string) error {
		err := validateParams(p, "Location")
		if err != nil {
			return err
		}
		bldr.SetLocation(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate CompanyName if provided
	if vldErr := setNillableStringField(p.CompanyName, func(v string) error {
		err := validateParams(p, "CompanyName")
		if err != nil {
			return err
		}
		bldr.SetCompanyName(v)
		return nil
	}); err != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Description if provided
	if vldErr := setNillableStringField(p.Description, func(program string) error {
		err := validateParams(p, "Description")
		if err != nil {
			return err
		}
		bldr.SetDescription(p.Description)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate JobTitle if provided
	if vldErr := setNillableStringField(p.JobTitle, func(v string) error {
		err := validateParams(p, "JobTitle")
		if err != nil {
			return err
		}
		bldr.SetJobTitle(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate PrimaryTechnology if provided
	if vldErr := setNillableJSONArrayField(p.PrimaryTechnologies, func(v []string) error {
		err := validateParams(p, "PrimaryTechnologies")
		if err != nil {
			return err
		}
		bldr.SetPrimaryTechnologies(p.PrimaryTechnologies)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Return all validation errors at once
	// this prevents the client from making several round trips to the server
	if collection.HasAny(vldErrs) {
		return nil, vldErrs
	}

	record, err = bldr.Save(dBContext)
	if err != nil {
		return nil, []error{err}
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

func (r *WorkExperienceRepository) GetWorkExperienceByTalentUUID(talentID int) (*ent.WorkExperience, error) {
	record, err := dBConn.WorkExperience.Query().
		Where(workexperience.And(
			workexperience.TalentID(talentID),
			workexperience.DeletedAtIsNil())).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	return record, nil
}
