package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/job"
	"github.com/10hourlabs/tentn/internal/database"
	"github.com/google/uuid"
	"github.com/gosimple/slug"

	_ "github.com/joho/godotenv/autoload"
)

type JobService struct {
	psqlClient   *ent.Client
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
		psqlClient:   client,
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
	// TODO: look into the performance of this query
	// I don't think it's a concerning but just so we
	// understand what query the ent orm generates
	// and see if we could make possible optimization
	job, err := js.psqlClient.Job.Query().
		Where(job.UUIDEQ(jobUUID)).
		Only(js.queryContext)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (js *JobService) DeleteJob(jobUUID uuid.UUID) error {
	// TODO: should we check if the job is
	// already updated before allowing UPDATE?
	// this will prevent unecessary updates
	_, err := js.psqlClient.Job.Update().
		Where(job.UUIDEQ(jobUUID)).
		SetDeletedAt(time.Now()).
		SetHiring(false).
		Save(js.queryContext)
	if err != nil {
		return err
	}
	return nil
}

func (js *JobService) UpdateJob(jobUUID uuid.UUID, j *ent.Job) (*ent.Job, error) {
	id, err := js.psqlClient.Job.Update().
		Where(job.UUIDEQ(jobUUID)).
		SetHiring(j.Hiring).
		SetTitle(j.Title).
		SetSummary(j.Summary).
		SetEmployment(j.Employment).
		SetCategory(j.Category).
		SetThumbnail(j.Thumbnail).
		SetWeHave(j.WeHave).
		SetRequirements(j.Requirements).
		SetYouHave(j.YouHave).
		Save(js.queryContext)
	if err != nil {
		return nil, err
	}
	// TODO: remove unecessary GET call
	j, err = js.psqlClient.Job.Get(js.queryContext, id)
	if err != nil {
		return nil, err
	}
	return j, nil
}
