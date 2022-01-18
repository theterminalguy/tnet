package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/education"
	"github.com/10hourlabs/tentn/ent/predicate"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/10hourlabs/tentn/util/date"
	"github.com/google/uuid"
)

type EducationQuerier interface {
	GetAllForTalent(talentID int) ([]*ent.Education, error)
	GetAll() ([]*ent.Education, error)
	GetByUUID(id uuid.UUID) (*ent.Education, error)
	Create(p EducationParams) (*ent.Education, error)
	Update(id uuid.UUID, p EducationParams) (*ent.Education, []error)
	DeleteByUUID(id uuid.UUID) error
}

type EducationRepository struct{}

type EducationParams struct {
	TalentUUID      uuid.UUID `json:"talent_uuid" validate:"required"`
	InstitutionName string    `json:"institution_name" validate:"required"`
	Location        string    `json:"location" validate:"required"`
	Degree          string    `json:"degree" validate:"required"`
	Program         string    `json:"program" validate:"required"`
	Overview        string    `json:"overview" validate:"required"`
	StartDate       string    `json:"start_date" validate:"datetime=2006-01-02T15:04:05Z07:00"`
	EndDate         string    `json:"end_date"`
}

func NewEducationRepository() *EducationRepository {
	return &EducationRepository{}
}

func (*EducationRepository) Filter(prd ...predicate.Education) ([]*ent.Education, error) {
	educations, err := dBConn.Education.Query().
		Where(prd...).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return educations, nil
}

func (*EducationRepository) GetAllForTalent(talentID int) ([]*ent.Education, error) {
	records, err := dBConn.Education.Query().
		Where(education.And(
			education.TalentID(talentID),
			education.DeletedAtIsNil())).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*EducationRepository) GetAll() ([]*ent.Education, error) {
	records, err := dBConn.Education.Query().
		Where(education.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*EducationRepository) GetByUUID(id uuid.UUID) (*ent.Education, error) {
	record, err := dBConn.Education.Query().
		Where(education.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*EducationRepository) Create(p EducationParams) (*ent.Education, error) {
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

	record, err := dBConn.Education.
		Create().
		SetTalentID(a.ID).
		SetDegree(p.Degree).
		SetInstitutionName(p.InstitutionName).
		SetLocation(p.Location).
		SetProgram(p.Program).
		SetOverview(p.Overview).
		SetStartDate(*sd).
		SetNillableEndDate(ed).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *EducationRepository) Update(id uuid.UUID, p EducationParams) (*ent.Education, []error) {
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

	// Set and Validate Degree if provided
	if vldErr := setNillableStringField(p.Degree, func(v string) error {
		err := validateParams(p, "Degree")
		if err != nil {
			return err
		}
		bldr.SetDegree(v)
		return nil
	}); err != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Program if provided
	if vldErr := setNillableStringField(p.Program, func(program string) error {
		err := validateParams(p, "Program")
		if err != nil {
			return err
		}
		bldr.SetProgram(p.Program)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Overview if provided
	if vldErr := setNillableStringField(p.Overview, func(v string) error {
		err := validateParams(p, "Overview")
		if err != nil {
			return err
		}
		bldr.SetOverview(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate InsitutionName if provided
	if vldErr := setNillableStringField(p.InstitutionName, func(v string) error {
		err := validateParams(p, "InstitutionName")
		if err != nil {
			return err
		}
		bldr.SetInstitutionName(p.InstitutionName)
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

func (r *EducationRepository) DeleteByUUID(id uuid.UUID) error {
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

func (r *EducationRepository) GetEducationByTalentUUID(talentID int) (*ent.Education, error) {
	record, err := dBConn.Education.Query().
		Where(education.And(
			education.TalentID(talentID),
			education.DeletedAtIsNil())).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	return record, nil
}
