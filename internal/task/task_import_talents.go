package task

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/google/uuid"
)

type ImportTalents struct {
	UserRepo      *repo.UserRepository
	TalentService *service.TalentService
}

func NewImportTalents() *ImportTalents {
	return &ImportTalents{
		UserRepo:      repo.NewUserRepository(),
		TalentService: service.NewTalentService(),
	}
}

func (t *ImportTalents) CreateFakeTalent() error {
	return nil
}

type UsertUserParams struct {
	ID                    uuid.UUID
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	Pronoun               string `json:"pronoun"`
	PreferredJobTitle     string `json:"preferred_job_title"`
	Email                 string `json:"email"`
	Phone                 string `json:"phone"`
	ProfessionSummary     string `json:"profession_summary"`
	City                  string `json:"city"`
	CountryCode           string `json:"country_code"`
	ProfessionalStartDate string `json:"professional_start_date"`
}

func (t *ImportTalents) Run(_ string) error {
	// read CSV file
	// for each row
	//   create user
	//   create talent

	f, err := os.Open("/home/theterminalguy/Downloads/talents.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	index := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(record[1])
		if index == 5 {
			break
		}
		index++
	}
	return nil
}
