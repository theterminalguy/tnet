package task

import (
	"fmt"
	"math/rand"

	"github.com/theterminalguy/tentn/ent/schema"
	"github.com/theterminalguy/tentn/ent/talent"
	repo "github.com/theterminalguy/tentn/internal/repository"
	"github.com/theterminalguy/tentn/internal/service"
	"github.com/theterminalguy/tentn/util/collection"
	"github.com/theterminalguy/tentn/util/date"
	faker "github.com/brianvoe/gofakeit/v6"
)

type TaskCreateFakeTalents struct {
	UserRepo      *repo.UserRepository
	TalentService *service.TalentService
}

func NewTaskCreateFakeTalents() *TaskCreateFakeTalents {
	return &TaskCreateFakeTalents{
		UserRepo:      repo.NewUserRepository(),
		TalentService: service.NewTalentService(),
	}
}

func (t *TaskCreateFakeTalents) CreateFakeTalent() error {
	fName := faker.FirstName()
	lName := faker.LastName()
	userParams := repo.UserParams{
		FirstName: fName,
		LastName:  lName,
		PhotoURL:  faker.ImageURL(128, 128),
		Email:     faker.Email(),
		Role:      "talent",
		Approved:  true,
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
		JobPreference: (func() talent.JobPreference {
			count := len(schema.JobPreferences())
			randomTalent := schema.JobPreferences()[rand.Intn(count)]
			return talent.JobPreference(randomTalent)
		})(),
		TimeZone: (func() string {
			return []string{"UTC", "WAT", "EET", "AET"}[rand.Intn(2)]
		})(),
		Available: faker.Bool(),
		State:     faker.State(),
	}
	_, err = t.TalentService.CreateProfile(user, talentParams)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskCreateFakeTalents) Run(_ string) error {
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
