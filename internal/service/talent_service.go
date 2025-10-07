package service

import (
	"github.com/theterminalguy/tentn/ent"
	"github.com/theterminalguy/tentn/ent/schema/userrole"
	"github.com/theterminalguy/tentn/internal/decorator"
	repo "github.com/theterminalguy/tentn/internal/repository"
)

type TalentService struct {
	UserRepo   *repo.UserRepository
	TalentRepo *repo.TalentRepository
}

func NewTalentService() *TalentService {
	return &TalentService{
		UserRepo:   repo.NewUserRepository(),
		TalentRepo: repo.NewTalentRepository(),
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

func (t *TalentService) CreateProfile(user *ent.User, p repo.TalentParams) (*decorator.TalentResponse, error) {
	if p.FirstName == "" {
		p.FirstName = user.FirstName
	}
	if p.LastName == "" {
		p.LastName = user.LastName
	}
	p.Email = user.Email
	p.UserID = user.ID
	p.Available = true
	a, err := t.TalentRepo.Create(p)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (t *TalentService) UpdateProfile(user *ent.User, p *repo.TalentParams) (*decorator.TalentResponse, []error) {
	p.Email = user.Email
	talent, err := t.TalentRepo.GetTalentByUserID(user.ID)
	if err != nil {
		return nil, []error{err}
	}
	a, vldErrs := t.TalentRepo.Update(talent.ID, *p)
	if vldErrs != nil {
		return nil, vldErrs
	}
	return a, nil
}
