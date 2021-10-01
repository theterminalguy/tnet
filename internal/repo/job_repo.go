package repo

import (
	"context"
	"os"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/job"
	"github.com/10hourlabs/tentn/internal/database"
	"github.com/google/uuid"
)

type JobRepository struct {
	db           *ent.Client
	queryContext context.Context
}

func NewJobRepository() *JobRepository {
	// TODO: see NewJobService
	client, err := database.NewPostgresClient(os.Getenv("TENTN_POSTGRES_DSN"))
	if err != nil {
		return nil
	}
	return &JobRepository{
		db:           client,
		queryContext: context.Background(),
	}
}

func (jr *JobRepository) GetAll() ([]*ent.Job, error) {
	jobs, err := jr.db.Job.Query().All(jr.queryContext)
	if err != nil {
		return []*ent.Job{}, err
	}
	return jobs, nil
}

func (jr *JobRepository) GetByUUID(jobUUID uuid.UUID) (*ent.Job, error) {
	job, err := jr.db.Job.Query().
		Where(job.UUIDEQ(jobUUID)).
		Only(jr.queryContext)
	if err != nil {
		return nil, err
	}
	return job, nil
}

// TODO: implement save
// func (jr *JobRepository) Save(j *ent.Job) (*ent.Job, error) {
// 	j.Update().SetCategory()
// }
