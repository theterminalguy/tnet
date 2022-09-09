package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/job"
	"github.com/10hourlabs/tentn/ent/predicate"
	"github.com/10hourlabs/tentn/internal/paginator"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/google/uuid"
)

type JobQuerier interface {
	GetAllForRecruiter(id uuid.UUID) ([]*ent.Job, error)
	GetAll(page string) (*paginator.OffsetPaginater, error)
	GetByID(id uuid.UUID) (*ent.Job, error)
	Create(p JobParams) (*ent.Job, error)
	Update(id uuid.UUID, p JobParams) (*ent.Job, []error)
	DeleteByID(id uuid.UUID) error
}

type JobRepository struct {
	TalentCollectionRepo TalentCollectionRepository
}

type JobParams struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`

	// TODO: automatically expire job
	Hiring     bool   `json:"hiring"`
	Title      string `json:"job_title" validate:"required"`
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
	TimeZone     string   `json:"timezone_id" validate:"required"`
	Location     string   `json:"location"`

	AtsJobID           string    `json:"ats_job_id" validate:"required"`
	TalentCollectionId uuid.UUID `json:"talent_collection_id"`

	/*
		TODO: Add support for slary range
		StartingSalary decimal.Decimal `json:"starting_salary" validate:"required"`
		Currency       string          `json:"currency" validate:"required"`
		EndingSalary   decimal.Decimal `json:"ending_salary" validate:"required"`
	*/
}

func NewJobRepository() *JobRepository {
	return &JobRepository{
		TalentCollectionRepo: *NewTalentCollectionRepository(),
	}
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

func (*JobRepository) GetAllForRecruiter(id uuid.UUID) ([]*ent.Job, error) {
	jobs, err := dBConn.Job.Query().
		Where(job.UserID(id),
			job.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (*JobRepository) GetAll(page string) (*paginator.OffsetPaginater, error) {
	pager, err := paginator.NewOffsetPaginater(page)
	if err != nil {
		return nil, err
	}
	jobs, err := dBConn.Job.Query().
		Limit(paginator.MaxResults).
		Offset(pager.GetOffset()).
		Where(job.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	// convert jobs to an interface array
	var jobList []interface{}
	for _, j := range jobs {
		jobList = append(jobList, j)
	}
	return pager.Paginate(jobList, len(jobList)), nil
}

func (*JobRepository) GetByID(id uuid.UUID) (*ent.Job, error) {
	j, err := dBConn.Job.Query().
		Where(job.ID(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if j.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return j, nil
}

func (jp *JobRepository) Create(p JobParams) (*ent.Job, error) {
	tx, err := dBConn.Tx(dBContext)
	if err != nil {
		return nil, fmt.Errorf("starting a transaction: %w", err)
	}
	err = ValidateParams(p)
	if err != nil {
		return nil, err
	}
	timeZoneName := TimeZoneRepo[p.TimeZone]
	if timeZoneName[1] == "" {
		return nil, errors.New("timezone not allowed")
	}

	// Create collection
	collectionParams := new(TalentCollectionParams)
	collectionParams.UserID = p.UserID
	collectionParams.Name = fmt.Sprintf("%s-%s", p.Title, p.AtsJobID)
	if err != nil {
		return nil, err
	}
	err = jp.TalentCollectionRepo.validateScopedUniquenessOfName(collectionParams.Name, p.UserID)
	if err != nil {
		return nil, err
	}
	// convert uuids to strings
	TalentIDs := make([]string, len(collectionParams.TalentIDS))
	for i, uuid := range collectionParams.TalentIDS {
		TalentIDs[i] = uuid.String()
	}
	record, err := tx.TalentCollection.
		Create().
		SetName(collectionParams.Name).
		SetUserID(p.UserID).
		SetTalentUuids(TalentIDs).
		Save(dBContext)
	if err != nil {
		return nil, rollback(tx, fmt.Errorf("failed creating the job collection: %w", err))
	}
	if err != nil {
		return nil, err
	}
	p.TalentCollectionId = record.ID

	jobUUID := uuid.New()
	j, err := tx.Job.
		Create().
		SetID(jobUUID).
		SetHiring(p.Hiring).
		SetTitle(p.Title).
		SetSummary(p.Summary).
		SetAtsJobID(p.AtsJobID).
		SetSlug(slugify(p.Title, jobUUID)).
		SetEmployment(job.Employment(p.Employment)).
		SetCategory(job.Category(p.Category)).
		SetThumbnail(p.Thumbnail).
		SetWeHave(p.WeHave).
		SetRequirements(p.Requirements).
		SetYouHave(p.YouHave).
		SetUserID(p.UserID).
		SetTalentCollectionID(p.TalentCollectionId).
		SetTimezone(timeZoneName[1]).
		Save(dBContext)
	if err != nil {
		return nil, rollback(tx, fmt.Errorf("failed creating the job: %w", err))
	}
	return j, tx.Commit()
}

func (r *JobRepository) Update(id uuid.UUID, p JobParams) (*ent.Job, []error) {
	record, err := r.GetByID(id)
	if err != nil {
		return nil, []error{err}
	}

	var vldErrs []error
	bldr := record.Update()

	// Set and Validate Hiring if provided
	if vldErr := setNillableBoolField(p.Hiring, func(v bool) error {
		err := ValidateParams(p, "Hiring")
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
		err := ValidateParams(p, "Title")
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
		err := ValidateParams(p, "Summary")
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
		err := ValidateParams(p, "Employment")
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
		err := ValidateParams(p, "Category")
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
		err := ValidateParams(p, "Thumbnail")
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
		err := ValidateParams(p, "WeHave")
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
		err := ValidateParams(p, "Requirements")
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
		err := ValidateParams(p, "YouHave")
		if err != nil {
			return err
		}
		bldr.SetYouHave(p.YouHave)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate TimeZone if provided
	if vldErr := setNillableStringField(string(p.TimeZone), func(v string) error {
		err := ValidateParams(p, "TimeZone")
		if err != nil {
			return err
		}
		timeZoneName := TimeZoneRepo[p.TimeZone]
		if timeZoneName[1] == "" {
			return errors.New("timezone not allowed")
		}
		bldr.SetTimezone(p.TimeZone)
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

func (r *JobRepository) DeleteByID(id uuid.UUID) error {
	record, err := r.GetByID(id)
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
