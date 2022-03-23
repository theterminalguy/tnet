package scope

import (
	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/job"
	"github.com/10hourlabs/tentn/ent/jobapplication"
	"github.com/10hourlabs/tentn/ent/talentcollection"
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

func (r *RecruiterScope) GetTalentCollections() ([]*ent.TalentCollection, error) {
	return r.Recruiter.QueryTalentCollections().All(repo.GetDBContext())
}

func (r *RecruiterScope) DeleteCollection(id uuid.UUID) error {
	record, err := r.Recruiter.QueryTalentCollections().Where(
		talentcollection.UserIDEQ(r.Recruiter.ID),
		talentcollection.UUIDEQ(id),
	).First(repo.GetDBContext())
	if err != nil {
		return err
	}
	tRepo := repo.NewTalentCollectionRepository()
	return tRepo.DeleteByUUID(record.UUID)
}
