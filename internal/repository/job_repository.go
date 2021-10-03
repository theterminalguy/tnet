package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/job"
	"github.com/google/uuid"
)

type JobRepository struct{}

type JobParams struct {
	Hiring       bool     `json:"hiring"`
	Title        string   `json:"title" validate:"required"`
	Summary      string   `json:"summary" validate:"required"`
	Employment   string   `json:"employment" validate:"required"`
	Category     string   `json:"category" validate:"required"`

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
		return nil, RecordNotFoundError
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

func (r *JobRepository) Update(id uuid.UUID, p JobParams) (*ent.Job, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	j, err := r.GetByUUID(id)
	if err != nil {
		return nil, err
	}
	_, err = dBConn.Job.Update().
		SetHiring(p.Hiring).
		SetTitle(p.Title).
		SetSummary(p.Summary).
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
	return j, nil
}

func (r *JobRepository) DeleteByUUID(id uuid.UUID) error {
	_, err := r.GetByUUID(id)
	if err != nil {
		return err
	}
	_, err = dBConn.Job.Update().
		SetDeletedAt(time.Now()).
		SetHiring(false).
		Save(dBContext)
	if err != nil {
		return err
	}
	return nil
}
