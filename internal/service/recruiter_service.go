package service

import (
	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	"github.com/10hourlabs/tentn/ent/slackappinstall"
	repo "github.com/10hourlabs/tentn/internal/repository"
)

type RecruiterService struct {
	UserRepo       *repo.UserRepository
	SlackAppRepo   *repo.SlackAppInstallRepository
	CollectionRepo *repo.TalentCollectionRepository
}

func NewRecruiterService() *RecruiterService {
	return &RecruiterService{
		UserRepo:     repo.NewUserRepository(),
		SlackAppRepo: repo.NewSlackAppInstallRepository(),
	}
}

func (rs *RecruiterService) InstallSlackApp(up repo.UserParams, sp repo.SlackAppInstallParams) (*ent.User, error) {
	// check if slack app already exists
	app, err := rs.SlackAppRepo.GetByTeamID(sp.TeamID)
	if app != nil {
		return app.Edges.User, nil
	}
	if err == repo.ErrRecordDeleted {
		app, err := rs.SlackAppRepo.GetDeletedInstallation(sp.TeamID)
		if err != nil {
			return nil, err
		}
		if err := rs.SlackAppRepo.UpdateFields(app, map[string]interface{}{
			slackappinstall.FieldDeletedAt:           nil,
			slackappinstall.FieldInstallCount:        app.InstallCount + 1,
			slackappinstall.FieldAccessToken:         sp.AccessToken,
			slackappinstall.FieldTeamName:            sp.TeamName,
			slackappinstall.FieldBotUserID:           sp.BotUserID,
			slackappinstall.FieldIsEnterpriseInstall: sp.IsEnterpriseInstall,
			slackappinstall.FieldScope:               sp.Scope,
		}); err != nil {
			return nil, err
		}
		return app.Edges.User, nil
	}
	if ent.IsNotFound(err) {
		up.Approved = true // Approve recruiter
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
		// create a new "Favorite" collection for the recruiter
		rs.CollectionRepo.Create(repo.TalentCollectionParams{
			UserID: recruiter.ID,
			Name:   "Favorite",
		})
		return recruiter, nil
	}
	return nil, err
}
