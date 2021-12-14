package service

import (
	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	repo "github.com/10hourlabs/tentn/internal/repository"
)

type TalentService struct {
	UserRepo *repo.UserRepository
}

func NewTalentService() *TalentService {
	return &TalentService{
		UserRepo: repo.NewUserRepository(),
	}
}

func (t *TalentService) RegisterTalent(up *repo.UserParams) (*ent.User, error) {
	ur := t.UserRepo
	talent, err := ur.GetByEmail(up.Email)
	if talent != nil {
		return talent, nil
	}
	if ent.IsNotFound(err) {
		up.Role = userrole.Talent
		up.Approved = true
		talent, err = ur.Create(*up)
		if err != nil {
			return nil, err
		}
		return talent, nil
	}
	return nil, err
}
