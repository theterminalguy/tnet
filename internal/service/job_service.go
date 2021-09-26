package service

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/job"
	"github.com/10hourlabs/tentn/internal/database"
	"github.com/google/uuid"
	"github.com/gosimple/slug"

	_ "github.com/joho/godotenv/autoload"
)

type JobService struct {
	psqlClient *ent.Client
	queryContext context.Context
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
		queryContext: context.Background(),
	}, nil
}

func (js *JobService) CreateJob(job *ent.Job) (*ent.Job, error) {
	jobUUID := uuid.New()
	jobSlug := slug.Make(fmt.Sprintf("%v %v", job.Title, jobUUID))
	job, err := js.psqlClient.Job.
		Create().
		SetUUID(jobUUID).
		SetHiring(job.Hiring).
		SetTitle(job.Title).
		SetSummary(job.Summary).
		SetSlug(jobSlug).
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
	jobs, err := js.psqlClient.Job.Query().All(js.queryContext)
	if err != nil {
		return nil, err
	}
	// TODO: remove logs
	log.Println("found jobs", jobs)
	return jobs, nil
}

func (js *JobService) GetJob(jobUUID uuid.UUID) (*ent.Job, error) {
	job, err := js.psqlClient.Job.Query().
		Where(job.UUIDEQ(jobUUID)).
		Only(js.queryContext)
	if err != nil {
		return nil, err
	}
	return job, nil
}
