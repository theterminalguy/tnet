package scope

import (
	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/job"
	"github.com/10hourlabs/tentn/ent/jobapplication"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/google/uuid"
)

type RecruiterScope struct {
	Recruiter *ent.User
}

func NewRecruiterScope(recruiter *ent.User) *RecruiterScope {
	return &RecruiterScope{
		Recruiter: recruiter,
	}
}

func (r *RecruiterScope) GetJobs() ([]*ent.Job, error) {
	return r.Recruiter.QueryJobs().
		WithApplications(func(jaq *ent.JobApplicationQuery) {
			jaq.WithTalent()
		}).All(repo.GetDBContext())
}

func (r *RecruiterScope) GetJobByUUID(uuid uuid.UUID) (*ent.Job, error) {
	return r.Recruiter.QueryJobs().
		Where(job.UUIDEQ(uuid)).
		WithApplications(func(jaq *ent.JobApplicationQuery) {
			jaq.WithTalent()
		}).
		First(repo.GetDBContext())
}

func (r *RecruiterScope) GetJobApplicationByUUID(uuid uuid.UUID) (*ent.JobApplication, error) {
	return r.Recruiter.QueryJobs().QueryApplications().WithTalent().WithJob().
		Where(jobapplication.UUIDEQ(uuid)).
		First(repo.GetDBContext())
}

func (r *RecruiterScope) UpdateJobApplication(uuid uuid.UUID, params repo.JobApplicationParams) (*ent.JobApplication, []error) {
	return repo.NewJobApplicationRepository().Update(uuid, params)
}
