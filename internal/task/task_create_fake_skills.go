package task

import (
	"errors"
	"math/rand"

	repo "github.com/theterminalguy/tentn/internal/repository"
	mathutil "github.com/theterminalguy/tentn/util/math"
	faker "github.com/brianvoe/gofakeit/v6"
)

type TaskCreateFakeSkill struct {
	TalentRepo *repo.TalentRepository
	SkillRepo  *repo.SkillRepository
}

func NewTaskCreateFakeSkill() *TaskCreateFakeSkill {
	return &TaskCreateFakeSkill{
		TalentRepo: repo.NewTalentRepository(),
		SkillRepo:  repo.NewSkillRepository(),
	}
}

func (c *TaskCreateFakeSkill) Run(_ string) error {
	talents, err := c.TalentRepo.GetAll()
	if err != nil {
		return err
	}
	if len(talents) == 0 {
		return errors.New("no talents found")
	}
	for _, talent := range talents {
		for i := 0; i < 3; i++ {
			sp := repo.SkillParams{
				TalentID:        talent.ID,
				YearsOfExperience: mathutil.RandomFloat32([]float32{}),
				Preferred: (func() bool {
					return []bool{true, false}[rand.Intn(2)]
				})(),
				Note: faker.Adjective(),
				Name: faker.ProgrammingLanguage(),
			}
			if _, err := c.SkillRepo.Create(sp); err != nil {
				return err
			}
		}
	}
	return nil
}
