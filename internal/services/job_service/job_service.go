package job_service

import (
	"context"
	"log"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/internal/database"
)

type Service struct {
	psqlClient *ent.Client
}

func NewJobService() (*Service, error) {
	client, err := database.NewPostgresClient()
	if err != nil {
		return nil, err
	}
	return &Service{
		psqlClient: client,
	}, nil
}

func (js *Service) CreateJob(job *ent.Job) (*ent.Job, error) {
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

func (js *Service) GetAllJobs() ([]*ent.Job, error) {
	defer js.psqlClient.Close()

	jobs, err := js.psqlClient.Job.Query().All(context.Background())
	if err != nil {
		return nil, err
	}
	// TODO: remove logs
	log.Println("found jobs", jobs)
	return jobs, nil
}
