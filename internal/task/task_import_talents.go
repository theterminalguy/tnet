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
	FirstName             string                      `json:"first_name"`
	LastName              string                      `json:"last_name"`
	Pronoun               string                      `json:"pronoun"`
	PreferredJobTitle     string                      `json:"preferred_job_title"`
	Email                 string                      `json:"email"`
	Phone                 string                      `json:"phone"`
	ProfessionalSummary   string                      `json:"professional_summary"`
	City                  string                      `json:"city"`
	CountryCode           string                      `json:"country_code"`
	ProfessionalStartDate string                      `json:"professional_start_date"`
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
	var uup UpsertUserParams
	err := json.Unmarshal(record, &uup)
	if err != nil {
		fmt.Println("error unmarshalling: ", err)
	}
	user := &repo.UserParams{
		ID:        uuid.New(),
		FirstName: uup.FirstName,
		LastName:  uup.LastName,
		Email:    uup.Email,
		PhotoURL: fmt.Sprintf(
			"https://ui-avatars.com/api/?name=%s+%s&background=random&size=64",
			uup.FirstName,
			uup.LastName,
		),
		Role:     "talent",
		Approved: true,
	}

	fmt.Println(user.FirstName + " " + user.LastName)
	fmt.Println("\tID: ", user.ID)
	fmt.Println("\tFirst Name: ", user.FirstName)
	fmt.Println("\tLast Name: ", user.LastName)
	fmt.Println("\tPhoto URL: ", user.PhotoURL)
	fmt.Println("\tEmail: ", user.Email)
	fmt.Println("\tRole: ", user.Role)
	fmt.Println("\tApproved: ", user.Approved)
	fmt.Print("==========================\n\n")

	talent := &repo.TalentParams{
		ID:                    uuid.New(),
		UserID:                user.ID,
		FirstName:             user.FirstName,
		LastName:              user.LastName,
		Email:                 user.Email,
		PreferredName:         user.FirstName,
		Pronoun:               uup.Pronoun,
		PreferredJobTitle:     uup.PreferredJobTitle,
		ProfessionalStartDate: uup.ProfessionalStartDate,
		Phone:                 uup.Phone,
		CountryCode:           uup.CountryCode,
		City:                  uup.City,
		JobPreference:         "flexible",
		TimeZone:              "GMT+1",
		Available:             true,
		State:                 uup.City,
		ProfessionalSummary:   fmt.Sprintf("An experienced %s", uup.PreferredJobTitle),
	}

	fmt.Println("\tTalent: ")
	fmt.Println("\t\tID: ", talent.ID)
	fmt.Println("\t\tUserID: ", talent.UserID)
	fmt.Println("\t\tFirst Name: ", talent.FirstName)
	fmt.Println("\t\tLast Name: ", talent.LastName)
	fmt.Println("\t\tEmail: ", talent.Email)
	fmt.Println("\t\tPreferred Name: ", talent.PreferredName)
	fmt.Println("\t\tPronoun: ", talent.Pronoun)
	fmt.Println("\t\tPreferred Job Title: ", talent.PreferredJobTitle)
	fmt.Println("\t\tPhone: ", talent.Phone)
	fmt.Println("\t\tCountry Code: ", talent.CountryCode)
	fmt.Println("\t\tCity: ", talent.City)
	fmt.Println("\t\tState: ", talent.State)
	fmt.Println("\t\tTimeZone: ", talent.TimeZone)
	fmt.Println("\t\tJob Preference: ", talent.JobPreference)
	fmt.Println("\t\tAvailable: ", talent.Available)
	fmt.Println("\t\tProfessional Summary: ", talent.ProfessionalSummary)

}
