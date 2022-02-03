package task

import (
	"fmt"
	"math/rand"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/10hourlabs/tentn/util/date"
	faker "github.com/brianvoe/gofakeit/v6"
)

type CreateFakeTalents struct {
	UserRepo      *repo.UserRepository
	TalentService *service.TalentService
}

func NewCreateFakeTalents() *CreateFakeTalents {
	return &CreateFakeTalents{
		UserRepo:      repo.NewUserRepository(),
		TalentService: service.NewTalentService(),
	}
}

func (t *CreateFakeTalents) CreateFakeTalent() error {
	fName := faker.FirstName()
	lName := faker.LastName()
	userParams := repo.UserParams{
		Name:     fName + " " + lName,
		Email:    faker.Email(),
		Role:     "talent",
		Approved: true,
	}
	user, err := t.UserRepo.Create(userParams)
	if err != nil {
		return err
	}
	talentParams := repo.TalentParams{
		UserID:        user.ID,
		FirstName:     fName,
		LastName:      lName,
		Email:         user.Email,
		PreferredName: faker.Username(),
		Pronoun: (func() string {
			return []string{"he", "she"}[rand.Intn(2)]
		})(),
		PreferredJobTitle:     faker.JobTitle(),
		ProfessionalStartDate: faker.Date().Format(date.ISOLayout),
		Phone:                 faker.Phone(),
		CountryCode:           faker.CountryAbr(),
		City:                  faker.City(),
		JobPreference: (func() []string {
			return [][]string{{"full_time", "part_time", "contract"}, {"full_time", "part_time"}, {"full_time"},
				{"contract"}, {"part_time"}, {"part_time", "contract"}, {"full_time", "contract"}}[rand.Intn(5)]
		})(),
		Available: faker.Bool(),
	}
	_, err = t.TalentService.CreateProfile(user, talentParams)
	if err != nil {
		return err
	}
	return nil
}

func (t *CreateFakeTalents) Run(_ string) error {
	var errs []error
	// TODO: make the max configurable
	// Also the inserts is not optimized, it works fine for now but should be improved
	for i := 0; i < 20; i++ {
		err := t.CreateFakeTalent()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if collection.HasAny(errs) {
		return fmt.Errorf("%d errors, %v", len(errs), errs)
	}
	return nil
}
