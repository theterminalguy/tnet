package task

import (
	"encoding/csv"
	"encoding/json"
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

type UpsertUserParams struct {
	ID                    uuid.UUID
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	Pronoun               string `json:"pronoun"`
	PreferredJobTitle     string `json:"preferred_job_title"`
	Email                 string `json:"email"`
	Phone                 string `json:"phone"`
	ProfessionalSummary   string `json:"professional_summary"`
	City                  string `json:"city"`
	CountryCode           string `json:"country_code"`
	ProfessionalStartDate string `json:"professional_start_date"`
	Talent                *repo.TalentParams
	Educations            []repo.EducationParams      `json:"educations"`
	WorkExperiences       []repo.WorkExperienceParams `json:"experiences"`
	Skills                []repo.SkillParams          `json:"skills"`
	Portfolios            []repo.PortfolioLinkParams  `json:"portfolios"`
}

func (t *ImportTalents) Run(_ string) error {
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
		payload := record[1]
		bulkUpsertTalent([]byte(payload))
		if index == 5 {
			break
		}
		index++
	}
	return nil
}

func bulkUpsertTalent(record []byte) {
	var newUser UpsertUserParams
	err := json.Unmarshal(record, &newUser)
	if err != nil {
		panic(err)
	}
	fmt.Println("ID: ", newUser.ID)
	fmt.Println("First Name: ", newUser.FirstName)
	fmt.Println("Last Name: ", newUser.LastName)
	fmt.Println("Pronoun: ", newUser.Pronoun)
	fmt.Println("Preferred Job Title: ", newUser.PreferredJobTitle)
	fmt.Println("Email: ", newUser.Email)
	fmt.Println("Phone: ", newUser.Phone)
	fmt.Println("Professional Summary", newUser.ProfessionalSummary)
	fmt.Println("City: ", newUser.City)
	fmt.Println("Country Code: ", newUser.CountryCode)
	fmt.Println("Start Date: ", newUser.ProfessionalStartDate)
	fmt.Println("==========================")
	fmt.Print("\n")
}
