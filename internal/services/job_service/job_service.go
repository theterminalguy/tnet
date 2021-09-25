package job_service

import (
	"context"
	"log"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/internal/database"
)

func CreateJob(job *ent.Job) (*ent.Job, error) {
	// TODO: move database to a service struct
	// so we don't have to call database in each method
	client, err := database.NewSQLite3InMemoryClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	job, err = client.Job.
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

func GetAllJobs() ([]*ent.Job, error) {
	client, err := database.NewSQLite3InMemoryClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	jobs, err := client.Job.Query().All(context.Background())
	if err != nil {
		return nil, err
	}
	// TODO: remove logs
	log.Println("found jobs", jobs)
	return jobs, nil
}
