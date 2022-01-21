package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/job"
	"github.com/10hourlabs/tentn/ent/predicate"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/google/uuid"
)

type JobQuerier interface {
	GetAll() ([]*ent.Job, error)
	GetByUUID(id uuid.UUID) (*ent.Job, error)
	Create(p JobParams) (*ent.Job, error)
	Update(id uuid.UUID, p JobParams) (*ent.Job, []error)
	DeleteByUUID(id uuid.UUID) error
}

type JobRepository struct{}

type JobParams struct {
	UserID     int    `json:"user_id" validate:"required"`
	Hiring     bool   `json:"hiring"`
	Title      string `json:"title" validate:"required"`
	Summary    string `json:"summary" validate:"required"`
	Employment string `json:"employment" validate:"required"`
	Category   string `json:"category" validate:"required"`

	// TODO thumbnail url should be validate against value provided
	// if we ever go public, not a concern for now
	// Go public in this context means if the api is made publicly available
	Thumbnail    string   `json:"thumbnail" validate:"required,url"`
	WeHave       []string `json:"we_have" validate:"required"`
	Requirements []string `json:"requirements" validate:"required"`
	YouHave      []string `json:"you_have" validate:"required"`
}

func NewJobRepository() *JobRepository {
	return &JobRepository{}
}

func (*JobRepository) Filter(prd ...predicate.Job) ([]*ent.Job, error) {
	jobs, err := dBConn.Job.Query().
		Where(prd...).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (*JobRepository) GetAll() ([]*ent.Job, error) {
	jobs, err := dBConn.Job.Query().
		Where(job.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (*JobRepository) GetByUUID(id uuid.UUID) (*ent.Job, error) {
	j, err := dBConn.Job.Query().
		Where(job.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if j.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return j, nil
}

func (*JobRepository) Create(p JobParams) (*ent.Job, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	jobUUID := uuid.New()
	j, err := dBConn.Job.
		Create().
		SetUUID(jobUUID).
		SetHiring(p.Hiring).
		SetTitle(p.Title).
		SetSummary(p.Summary).
		SetSlug(slugify(p.Title, jobUUID)).
		SetEmployment(job.Employment(p.Employment)).
		SetCategory(job.Category(p.Category)).
		SetThumbnail(p.Thumbnail).
		SetWeHave(p.WeHave).
		SetRequirements(p.Requirements).
		SetYouHave(p.YouHave).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return j, err
}

func (r *JobRepository) Update(id uuid.UUID, p JobParams) (*ent.Job, []error) {
	record, err := r.GetByUUID(id)
	if err != nil {
		return nil, []error{err}
	}

	var vldErrs []error
	bldr := record.Update()

	// Set and Validate Hiring if provided
	if vldErr := setNillableBoolField(p.Hiring, func(v bool) error {
		err := validateParams(p, "Hiring")
		if err != nil {
			return err
		}
		bldr.SetHiring(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Title if provided
	if vldErr := setNillableStringField(p.Title, func(v string) error {
		err := validateParams(p, "Title")
		if err != nil {
			return err
		}
		bldr.SetTitle(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Summary if provided
	if vldErr := setNillableStringField(p.Summary, func(v string) error {
		err := validateParams(p, "Summary")
		if err != nil {
			return err
		}
		bldr.SetSummary(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Employment if provided
	if vldErr := setNillableStringField(p.Employment, func(v string) error {
		err := validateParams(p, "Employment")
		if err != nil {
			return err
		}
		bldr.SetEmployment(job.Employment(v))
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Category if provided
	if vldErr := setNillableStringField(p.Category, func(v string) error {
		err := validateParams(p, "Category")
		if err != nil {
			return err
		}
		bldr.SetCategory(job.Category(v))
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Thumbnail if provided
	if vldErr := setNillableStringField(p.Thumbnail, func(v string) error {
		err := validateParams(p, "Thumbnail")
		if err != nil {
			return err
		}
		bldr.SetThumbnail(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate WeHave if provided
	if vldErr := setNillableJSONArrayField(p.WeHave, func(v []string) error {
		err := validateParams(p, "WeHave")
		if err != nil {
			return err
		}
		bldr.SetWeHave(p.WeHave)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Requirements if provided
	if vldErr := setNillableJSONArrayField(p.Requirements, func(v []string) error {
		err := validateParams(p, "Requirements")
		if err != nil {
			return err
		}
		bldr.SetRequirements(p.Requirements)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate YouHave if provided
	if vldErr := setNillableJSONArrayField(p.YouHave, func(v []string) error {
		err := validateParams(p, "YouHave")
		if err != nil {
			return err
		}
		bldr.SetYouHave(p.YouHave)
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

func (r *JobRepository) DeleteByUUID(id uuid.UUID) error {
	record, err := r.GetByUUID(id)
	if err != nil {
		return err
	}
	_, err = record.Update().
		SetDeletedAt(time.Now()).
		SetHiring(false).
		Save(dBContext)
	if err != nil {
		return err
	}
	return nil
}
