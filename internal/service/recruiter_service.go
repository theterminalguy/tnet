package service

import (
	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	repo "github.com/10hourlabs/tentn/internal/repository"
)

type RecruiterService struct {
	UserRepo     *repo.UserRepository
	SlackAppRepo *repo.SlackAppInstallRepository
}

func NewRecruiterService() *RecruiterService {
	return &RecruiterService{
		UserRepo:     repo.NewUserRepository(),
		SlackAppRepo: repo.NewSlackAppInstallRepository(),
	}
}

func (rs *RecruiterService) InstallSlackApp(up repo.UserParams, sp repo.SlackAppInstallParams) (*ent.User, error) {
	record, err := rs.UserRepo.GetByEmail(up.Email)
	if record != nil {
		return record, nil
	}
	if ent.IsNotFound(err) {
		up.Approved = false
		up.Role = userrole.Recruiter
		recruiter, err := rs.UserRepo.Create(up)
		if err != nil {
			return nil, err
		}
		// Create the slack app install record
		// Set the user id to the recruiter's id
		sp.UserID = recruiter.ID
		if _, err := rs.SlackAppRepo.Create(sp); err != nil {
			return nil, err
		}
		return recruiter, nil
	}
	return nil, err
}
