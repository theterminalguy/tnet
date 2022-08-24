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

func (r *RecruiterScope) GetID() uuid.UUID {
	return r.Recruiter.ID
}

func (r *RecruiterScope) GetJobs() ([]*ent.Job, error) {
	return r.Recruiter.QueryJobs().
		WithApplications(func(jaq *ent.JobApplicationQuery) {
			jaq.WithTalent()
		}).All(repo.GetDBContext())
}

func (r *RecruiterScope) GetJobByID(id uuid.UUID) (*ent.Job, error) {
	return r.Recruiter.QueryJobs().
		Where(job.ID(id)).
		WithJobPayments().
		WithApplications(func(jaq *ent.JobApplicationQuery) {
			jaq.WithTalent()
		}).
		First(repo.GetDBContext())
}

func (r *RecruiterScope) GetJobApplicationByID(id uuid.UUID) (*ent.JobApplication, error) {
	return r.Recruiter.QueryJobs().QueryApplications().WithTalent().WithJob().
		Where(jobapplication.ID(id)).
		First(repo.GetDBContext())
}

func (r *RecruiterScope) GetTalentCollections() ([]*ent.TalentCollection, error) {
	return r.Recruiter.QueryTalentCollections().All(repo.GetDBContext())
}

func (r *RecruiterScope) GetTalentCollectionByName(name string) (*ent.TalentCollection, error) {
	return r.Recruiter.QueryTalentCollections().
		Where(talentcollection.Name(name)).
		First(repo.GetDBContext())
}

func (r *RecruiterScope) GetTalentCollectionByID(id uuid.UUID) (*ent.TalentCollection, error) {
	return r.Recruiter.QueryTalentCollections().
		Where(
			talentcollection.And(
				talentcollection.UserIDEQ(r.Recruiter.ID),
				talentcollection.ID(id),
			)).First(repo.GetDBContext())
}

func (r *RecruiterScope) DeleteCollection(id uuid.UUID) error {
	record, err := r.Recruiter.QueryTalentCollections().Where(
		talentcollection.UserIDEQ(r.Recruiter.ID),
		talentcollection.ID(id),
	).First(repo.GetDBContext())
	if err != nil {
		return err
	}
	tRepo := repo.NewTalentCollectionRepository()
	return tRepo.DeleteByID(record.ID)
}

func (r *RecruiterScope) DeleteTalentsFromCollection(collectionID uuid.UUID, talentIDs []uuid.UUID) (*ent.TalentCollection, error) {
	talentCollection, err := r.Recruiter.QueryTalentCollections().Where(
		talentcollection.UserIDEQ(r.Recruiter.ID),
		talentcollection.ID(collectionID),
	).First(repo.GetDBContext())
	if err != nil {
		return nil, err
	}
	tRepo := repo.NewTalentCollectionRepository()
	return tRepo.RemoveTalents(talentCollection, talentIDs)
}

func (r *RecruiterScope) UpdateTalentCollection(collectionID uuid.UUID, params repo.TalentCollectionParams) (*ent.TalentCollection, error) {
	talentCollection, err := r.Recruiter.QueryTalentCollections().Where(
		talentcollection.UserIDEQ(r.Recruiter.ID),
		talentcollection.ID(collectionID),
	).First(repo.GetDBContext())
	if err != nil {
		return nil, err
	}
	tRepo := repo.NewTalentCollectionRepository()
	return tRepo.Update(talentCollection.ID, params)
}
