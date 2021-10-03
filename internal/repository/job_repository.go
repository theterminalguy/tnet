package repository

import (
	"log"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/job"
	"github.com/google/uuid"
)

type JobRepository struct{}

type JobParams struct {
	Hiring       bool     `json:"hiring"`
	Title        string   `json:"title"`
	Slug         string   `json:"slug"`
	Summary      string   `json:"summary"`
	Employment   string   `json:"employment"`
	Category     string   `json:"category"`
	Thumbnail    string   `json:"thumbnail"`
	WeHave       []string `json:"we_have"`
	Requirements []string `json:"requirements"`
	YouHave      []string `json:"you_have"`
}

func NewJobRepository() *JobRepository {
	return &JobRepository{}
}

func (*JobRepository) GetAll() ([]*ent.Job, error) {
	jobs, err := dBConn.Job.Query().
		Where(job.DeletedAtNotNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	// TODO: remove logs
	log.Println("found jobs", jobs)
	return jobs, nil
}

func (*JobRepository) GetByUUID(jobUUID uuid.UUID) (*ent.Job, error) {
	j, err := dBConn.Job.Query().
		Where(job.UUIDEQ(jobUUID)).
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
	// TODO: remove logs
	log.Println("job was created: ", j)
	return j, err
}

func (r *JobRepository) Update(jobUUID uuid.UUID, p JobParams) (*ent.Job, error) {
	j, err := r.GetByUUID(jobUUID)
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
	// TODO: potential for bug
	// does j gets updated after
	// a database update? :thinking_face:
	return j, nil
}

func (r *JobRepository) DeleteByUUID(id uuid.UUID) error {
	j, err := r.GetByUUID(id)
	if err != nil {
		return err
	}
	if j.Hiring == false {
		return RecordNotFoundError
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
