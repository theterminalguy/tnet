package job_service

import (
	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/internal/database"
)

func CreateJob(job *ent.Job) (*ent.Job, error) {
	client, err := database.Client()
	defer client.Close()
	if err != nil {
		return nil, err
	}
	return &ent.Job{}, nil
}
