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

type UpsertUserParams struct {
	ID                    uuid.UUID
	FirstName             string                       `json:"first_name"`
	LastName              string                       `json:"last_name"`
	Pronoun               string                       `json:"pronoun"`
	PreferredJobTitle     string                       `json:"preferred_job_title"`
	Email                 string                       `json:"email"`
	Phone                 string                       `json:"phone"`
	ProfessionalSummary   string                       `json:"professional_summary"`
	City                  string                       `json:"city"`
	CountryCode           string                       `json:"country_code"`
	ProfessionalStartDate string                       `json:"professional_start_date"`
	Educations            []*repo.EducationParams      `json:"educations"`
	WorkExperiences       []*repo.WorkExperienceParams `json:"experiences"`
	Skills                []*repo.SkillParams          `json:"skills"`
	Portfolios            []*repo.PortfolioLinkParams  `json:"portfolios"`
}

func (t *ImportTalents) Run(_ string) error {
	f, err := os.Open("/home/theterminalguy/Downloads/talents.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	var tds []*TalentData
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
		payload := []byte(record[1])
		tds = append(tds, extractTalentData(payload))
		if index == 20 {
			break
		}
		index++
	}
	fmt.Println(tds)
	return nil
}

type TalentData struct {
	User            *repo.UserParams
	Talent          *repo.TalentParams
	Skills          []*repo.SkillParams
	WorkExperiences []*repo.WorkExperienceParams
	Educations      []*repo.EducationParams
	Portfolios      []*repo.PortfolioLinkParams
}

func extractTalentData(record []byte) *TalentData {
	var uup UpsertUserParams
	err := json.Unmarshal(record, &uup)
	if err != nil {
		fmt.Println("error unmarshalling: ", err)
	}
	user := &repo.UserParams{
		ID:        uuid.New(),
		FirstName: uup.FirstName,
		LastName:  uup.LastName,
		Email:     uup.Email,
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
	fmt.Print("\t==========================\n\n")

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

	fmt.Println("\tSkills: ")
	for _, skill := range uup.Skills {
		// update skill IDs
		skill.ID = uuid.New()
		skill.TalentID = talent.ID
		skill.Preferred = true
		fmt.Println("\t\tID: ", skill.ID)
		fmt.Println("\t\tTalentID: ", skill.TalentID)
		fmt.Println("\t\tYears Of Experience: ", skill.YearsOfExperience)
		fmt.Println("\t\tPreferred: ", skill.Preferred)
		fmt.Println("\t\tNote: ", skill.Note)
		fmt.Println("\t\tName: ", skill.Name)
		fmt.Print("\t\t==========================\n\n")
	}

	fmt.Println("\tWork Experiences: ")
	for _, we := range uup.WorkExperiences {
		// update work experience IDs
		we.ID = uuid.New()
		we.TalentID = talent.ID
		fmt.Println("\t\tID: ", we.ID)
		fmt.Println("\t\tTalentID: ", we.TalentID)
		fmt.Println("\t\tCompany: ", we.CompanyName)
		fmt.Println("\t\tJob Title: ", we.JobTitle)
		fmt.Println("\t\tPrimary Technologies: ", we.PrimaryTechnologies)
		fmt.Println("\t\tDescription: ", we.Description)
		fmt.Println("\t\tStartDate: ", we.StartDate)
		fmt.Println("\t\tEndDate: ", we.EndDate)
		fmt.Print("\t\t==========================\n\n")
	}

	fmt.Println("\tEducations: ")
	for _, ed := range uup.Educations {
		// update education IDs
		ed.ID = uuid.New()
		ed.TalentID = talent.ID
		fmt.Println("\t\tID: ", ed.ID)
		fmt.Println("\t\tTalentID: ", ed.TalentID)
		fmt.Println("\t\tInstitution Name: ", ed.InstitutionName)
		fmt.Println("\t\tLocation: ", ed.Location)
		fmt.Println("\t\tDegree: ", ed.Degree)
		fmt.Println("\t\tProgram: ", ed.Program)
		fmt.Println("\t\tOverview: ", ed.Overview)
		fmt.Println("\t\tStartDate: ", ed.StartDate)
		fmt.Println("\t\tEndDate: ", ed.EndDate)
	}

	fmt.Println("\tPortfolios: ")
	for _, p := range uup.Portfolios {
		// update portfolio IDs
		p.ID = uuid.New()
		p.TalentID = talent.ID
		fmt.Println("\t\tID: ", p.ID)
		fmt.Println("\t\tTalentID: ", p.TalentID)
		fmt.Println("\t\tName: ", p.Name)
		fmt.Println("\t\tURL: ", p.URL)
	}

	fmt.Println("")

	return &TalentData{
		User:            user,
		Talent:          talent,
		Skills:          uup.Skills,
		Educations:      uup.Educations,
		WorkExperiences: uup.WorkExperiences,
		Portfolios:      uup.Portfolios,
	}
}
