package service

import (
	"errors"
	"os"

	"github.com/10hourlabs/tenlog"
	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/job"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/repository/scope"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/google/uuid"
)

type JobApplicationService struct {
	JobApplicationRepository *repo.JobApplicationRepository
	JobRepository            *repo.JobRepository
	TalentRepository         *repo.TalentRepository
	PortfolioLinkRepository  *repo.PortfolioLinkRepository
	EducationRepository      *repo.EducationRepository
	WorkExperienceRepository *repo.WorkExperienceRepository
	EmailTemplateRepository  *repo.EmailTemplateRepository
}

func NewJobApplicationService() *JobApplicationService {
	return &JobApplicationService{
		JobApplicationRepository: repo.NewJobApplicationRepository(),
		JobRepository:            repo.NewJobRepository(),
		TalentRepository:         repo.NewTalentRepository(),
		PortfolioLinkRepository:  repo.NewPortfolioLinkRepository(),
		EducationRepository:      repo.NewEducationRepository(),
		WorkExperienceRepository: repo.NewWorkExperienceRepository(),
		EmailTemplateRepository:  repo.NewEmailTemplateRepository(),
	}
}

func (*JobApplicationService) Apply(jobUUID, TalentID uuid.UUID) {
	// TODO
	// user must have a linkedin profile
	// user must have a GitHub profile for Enginering role
}

func (j *JobApplicationService) Validateq(jb ent.Job, pfLinks []*ent.PortfolioLink) error {

	// TODO: This maybe improved upon using ent edge query
	// example can be found here -- https://entgo.io/docs/predicates/#edge-predicates
	// Jira Ticket -- https://10hourlabs.atlassian.net/browse/MP-1

	//check if Talent has linkedin (required for every job category)
	hasLinkedIn := j.ContainsLinkedIn(pfLinks)

	if hasLinkedIn {
		switch jb.Category {

		case job.CategoryEngineering:
			return j.ContainsGithub(pfLinks)
		case job.CategoryProductDesign:
			return j.ContainsProductLinks(pfLinks)
		default:
			return nil
		}
	}
	return errors.New("no linkedin profile for talent")
}

//ContainsGithub checks to see if the array of given links contain Github
func (*JobApplicationService) ContainsGithub(pfLinks []*ent.PortfolioLink) error {
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
func (*JobApplicationService) ContainsProductLinks(pfLinks []*ent.PortfolioLink) error {
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
func (*JobApplicationService) ContainsLinkedIn(pfLinks []*ent.PortfolioLink) bool {
	var arr = make([]string, len(pfLinks))
	for _, val := range pfLinks {
		arr = append(arr, val.Name)
	}
	return contains(arr, "Linkedin")
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func (j *JobApplicationService) Create(p repo.JobApplicationParams) (*ent.JobApplication, error) {
	err := j.Validate(p.TalentID, p.JobUUID)
	if err != nil {
		return nil, err
	}
	record, err := j.JobApplicationRepository.Create(p)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (j *JobApplicationService) Validate(TalentID, jobUUID uuid.UUID) error {
	// check if talent exists
	talent, err := j.TalentRepository.GetByID(TalentID)
	if err != nil {
		return err
	}
	// check if job exists
	job, err := j.JobRepository.GetByID(jobUUID)
	if err != nil {
		return err
	}
	//check if user already applied for job
	err = j.JobApplicationRepository.AlreadyApplied(job.ID, talent.ID)
	if err != nil {
		return err
	}
	// check if user has education
	if !collection.HasAny(talent.Edges.Educations) {
		return errors.New("education not found")
	}
	// check if user has work experience
	_, err = j.WorkExperienceRepository.GetWorkExperienceByTalentID(talent.ID)
	if err != nil {
		return err
	}
	// check if user has preferred portfolio link
	if !collection.HasAny(talent.Edges.Portfoliolinks) {
		return errors.New("portfoliolink(s) not found")
	}

	err = j.Validateq(*job, talent.Edges.Portfoliolinks)
	if err != nil {
		return err
	}
	return nil
}

func (j *JobApplicationService) UpdateStatus(user *scope.RecruiterScope, id uuid.UUID, params repo.JobApplicationParams) (*ent.JobApplication, []error) {
	//update the status
	record, err := user.GetJobApplicationByID(id)
	if err != nil {
		return nil, []error{err}
	}
	jobApplication, queryErr := j.JobApplicationRepository.Update(record.ID, params)
	if err != nil {
		return nil, queryErr
	}
	// get the email the recruiter has saved for given status
	rctEmail, emailErr := j.EmailTemplateRepository.GetEmailForRecruiter(user.Recruiter.ID, params.Status)
	if emailErr != nil {
		return nil, []error{emailErr}
	}
	// call email sending api in the background
	go func() {
		_, err = SendMail(
			rctEmail.From,
			rctEmail.Subject,
			rctEmail.Body,
			record.Edges.Talent.Email,
			os.Getenv("DOMAIN"),
			os.Getenv("EMAIL_PRIVATE_API_KEY"),
			rctEmail.Cc,
			rctEmail.Bcc,
		)
		if err != nil {
			tenlog.Error(err)
		}
	}()
	return jobApplication, nil
}

func SendMail(args ...interface{}) (string, error) {
	return "", nil
}
