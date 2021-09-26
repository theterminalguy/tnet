package service

import (
	"context"
	"log"
	"os"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/internal/database"
)

type JobService struct {
	psqlClient *ent.Client
}

func NewJobService() (*JobService, error) {
	// TODO: consider moving database initialization
	// to a Service type that is then embedded into each individual service
	// the init() function of the service can setup all dependencies first
	client, err := database.NewPostgresClient(os.Getenv("TENTN_POSTGRES_DSN"))
	if err != nil {
		return nil, err
	}
	return &JobService{
		psqlClient: client,
	}, nil
}

func (js *JobService) CreateJob(job *ent.Job) (*ent.Job, error) {
	defer js.psqlClient.Close()

	job, err := js.psqlClient.Job.
		Create().
		SetHiring(job.Hiring).
		SetTitle(job.Title).
		SetSummary(job.Summary).
		// TODO slug should be automatically built
		SetSlug(job.Slug).
		SetEmployment(job.Employment).
		SetCategory(job.Category).
		SetThumbnail(job.Thumbnail).
		SetWeHave(job.WeHave).
		SetRequirements(job.Requirements).
		SetYouHave(job.YouHave).
		Save(context.Background())
	if err != nil {
		return nil, err
	}
	// TODO: remove logs
	log.Println("job was created: ", job)
	return job, nil
}

func (js *JobService) GetAllJobs() ([]*ent.Job, error) {
	defer js.psqlClient.Close()

	jobs, err := js.psqlClient.Job.Query().All(context.Background())
	if err != nil {
		return nil, err
	}
	// TODO: remove logs
	log.Println("found jobs", jobs)
	return jobs, nil
}
