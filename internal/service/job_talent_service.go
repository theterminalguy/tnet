package service

import (
	"errors"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/job"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/google/uuid"
)

type JobTalentService struct {
	JobTalentRepository     *repo.JobTalentRepository
	JobRepository           *repo.JobRepository
	TalentRepository        *repo.TalentRepository
	PortfolioLinkRepository *repo.PortfolioLinkRepository
}

func NewJobTalentService() *JobTalentService {
	return &JobTalentService{
		JobTalentRepository:     repo.NewJobTalentRepository(),
		JobRepository:           repo.NewJobRepository(),
		TalentRepository:        repo.NewTalentRepository(),
		PortfolioLinkRepository: repo.NewPortfolioLinkRepository(),
	}
}

func (*JobTalentService) Apply(jobUUID, TalentUUID uuid.UUID) {
	// TODO
	// user must have a linkedin profile
	// user must have a GitHub profile for Enginering role
}

func (j *JobTalentService) Validate(TalentUUID, jobUUID uuid.UUID) error {

	// TODO: This maybe improved upon using ent edge query
	// example can be found here -- https://entgo.io/docs/predicates/#edge-predicates
	// Jira Ticket -- https://10hourlabs.atlassian.net/browse/MP-1

	// get Talent
	Talent, err := j.TalentRepository.GetByUUID(TalentUUID)
	if err != nil {
		return err
	}

	//check if Talent has linkedin (required for every job category)
	hasLinkedIn := j.ContainsLinkedIn(Talent.Edges.Portfoliolinks)

	if hasLinkedIn {
		// get the job
		jb, err := j.JobRepository.GetByUUID(jobUUID)
		if err != nil {
			return err
		}

		switch jb.Category {

		case job.CategoryEngineering:
			return j.ContainsGithub(Talent.Edges.Portfoliolinks)
		case job.CategoryProductDesign:
			return j.ContainsProductLinks(Talent.Edges.Portfoliolinks)
		default:
			return nil

		}

	}

	return errors.New("no linkedin profile for candidate")

}

//ContainsGithub checks to see if the array of given links contain Github
func (*JobTalentService) ContainsGithub(pfLinks []*ent.PortfolioLink) error {
	var arr = make([]string, len(pfLinks))
	for _, val := range pfLinks {
		arr = append(arr, val.Name)
	}
	if contains(arr, "Github") {
		return nil
	}

	return errors.New("need a github profile for this role")
}

//ContainsProductLinks checks to see if the array of given links contain either Dribble or Behance
func (*JobTalentService) ContainsProductLinks(pfLinks []*ent.PortfolioLink) error {
	var arr = make([]string, len(pfLinks))
	for _, val := range pfLinks {
		arr = append(arr, val.Name)
	}
	if contains(arr, "Dribble") || contains(arr, "Behance") {
		return nil
	}

	return errors.New("need a dribble or behance profile for this role")
}

//ContainsLinkedin checks to see if the array of given links contain LinkedIn
func (*JobTalentService) ContainsLinkedIn(pfLinks []*ent.PortfolioLink) bool {
	var arr = make([]string, len(pfLinks))
	for _, val := range pfLinks {
		arr = append(arr, val.Name)
	}
	return contains(arr, "LinkedIn")
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
