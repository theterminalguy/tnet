package task

import (
	"errors"
	"fmt"
	"math/rand"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/util/collection"
)

type CreateFakeJob struct {
	JobRepo *repo.JobRepository
}

func NewCreateFakeJob() *CreateFakeJob {
	return &CreateFakeJob{
		JobRepo: repo.NewJobRepository(),
	}
}

func (j *CreateFakeJob) CreateFakeJob(userID int) error {
	jobParams := repo.JobParams{
		Hiring: true,
		Title: (func() string {
			return []string{"Front-End Developer", "Back-End Developer"}[rand.Intn(2)]
		})(),
		Summary: (func() string {
			return []string{
				"Experienced with Unreal, Unity and Custom game engines? Join our team today",
				"Experienced with Web Development, using Golang, React, and NextJS",
			}[rand.Intn(2)]
		})(),
		Employment: (func() string {
			return []string{"full_time", "part_time", "contract"}[rand.Intn(2)]
		})(),
		Category: (func() string {
			return []string{"engineering"}[rand.Intn(1)]
		})(),
		Thumbnail: (func() string {
			return []string{"https://picture.com/assets/images/e.png", "https://picture.com/assets/images/f.png"}[rand.Intn(2)]
		})(),
		WeHave: (func() []string {
			return [][]string{{"Casual and diverse workplace", "wifi"}}[rand.Intn(1)]
		})(),
		Requirements: (func() []string {
			return [][]string{{"Solving, and course correcting"}}[rand.Intn(1)]
		})(),
		YouHave: (func() []string {
			return [][]string{{
				"At least 4+ years of professional work experience as a Software Developer",
				"Understanding of OOP principles and practices",
				"Experience working with relational databases MySQL and PostgreSQL",
				"Version control and Git workflow",
				"You enjoy Front-End development",
				"Understand Containers, Infrastructure as Code and have worked with various cloud providers",
				"You enjoy writing tests",
			}}[rand.Intn(1)]
		})(),
		UserID: userID,
	}
	_, err := j.JobRepo.Create(jobParams)
	if err != nil {
		return err
	}
	return nil
}

func (j *CreateFakeJob) Run(_ string) error {
	var errs []error
	// TODO: make the max configurable
	// Also the inserts is not optimized, it works fine for now but should be improved
	jobs, err := j.JobRepo.GetAll()
	if err != nil {
		return err
	}
	if len(jobs) == 0 {
		return errors.New("no jobs found")
	}
	for _, job := range jobs {
		for i := 0; i < 3; i++ {
			err := j.CreateFakeJob(job.ID)
			if err != nil {
				errs = append(errs, err)
			}
			if collection.HasAny(errs) {
				return fmt.Errorf("%d errors, %v", len(errs), errs)
			}
		}
	}
	return nil
}
