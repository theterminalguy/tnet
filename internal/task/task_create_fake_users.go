package task

import (
	"fmt"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/util/collection"
	faker "github.com/brianvoe/gofakeit/v6"
)

type CreateFakeUsers struct {
	UserRepo *repo.UserRepository
}

func NewCreateFakeUsers() *CreateFakeTalents {
	return &CreateFakeTalents{
		UserRepo:      repo.NewUserRepository(),
		TalentService: service.NewTalentService(),
	}
}

func (t *CreateFakeUsers) CreateFakeUsers() error {
	fName := faker.FirstName()
	lName := faker.LastName()
	userParams := repo.UserParams{
		Name:     fName + " " + lName,
		Email:    faker.Email(),
		Role:     "talent",
		Approved: true,
	}
	_, err := t.UserRepo.Create(userParams)
	if err != nil {
		return err
	}
	return nil
}

func (t *CreateFakeUsers) Run(_ string) error {
	var errs []error
	// TODO: make the max configurable
	// Also the inserts is not optimized, it works fine for now but should be improved
	for i := 0; i < 20; i++ {
		err := t.CreateFakeUsers()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if collection.HasAny(errs) {
		return fmt.Errorf("%d errors, %v", len(errs), errs)
	}
	return nil
}
