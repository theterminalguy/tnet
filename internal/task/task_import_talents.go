package task

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/theterminalguy/tenlog"
	"github.com/theterminalguy/tnet/ent/schema/userrole"
	repo "github.com/theterminalguy/tnet/internal/repository"
	"github.com/theterminalguy/tnet/util/osutil"
	"github.com/theterminalguy/tnet/util/photo"
)

type TaskImportTalents struct {
	UserRepo      *repo.UserRepository
	TalentRepo    *repo.TalentRepository
	SkillsRepo    *repo.SkillRepository
	EduRepo       *repo.EducationRepository
	WorkRepo      *repo.WorkExperienceRepository
	PortfolioRepo *repo.PortfolioLinkRepository
}

func NewTaskImportTalents() *TaskImportTalents {
	return &TaskImportTalents{
		UserRepo:      repo.NewUserRepository(),
		TalentRepo:    repo.NewTalentRepository(),
		SkillsRepo:    repo.NewSkillRepository(),
		EduRepo:       repo.NewEducationRepository(),
		WorkRepo:      repo.NewWorkExperienceRepository(),
		PortfolioRepo: repo.NewPortfolioLinkRepository(),
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

func (t *TaskImportTalents) Run(_ string) error {
	manyUsers := make([]*repo.UserParams, 0)
	manyTalents := make([]*repo.TalentParams, 0)
	manySkills := make([]*repo.SkillParams, 0)
	manyPortfolios := make([]*repo.PortfolioLinkParams, 0)
	manyEducations := make([]*repo.EducationParams, 0)
	manyWorkExperiences := make([]*repo.WorkExperienceParams, 0)
	buildSearchIndex := make([]SearchIndexSchema, 0)

	var csvReader *csv.Reader
	resp, err := osutil.ReadCSVFromURL("https://storage.googleapis.com/tentn-bucket/temp/talents.csv")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	csvReader = csv.NewReader(resp.Body)
	emailCache := make(map[string]bool)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		payload := []byte(record[1])
		var uup UpsertUserParams
		err = json.Unmarshal(payload, &uup)
		if err != nil {
			fmt.Println("error unmarshalling: ", err)
			return err
		}
		usersEmail := strings.ToLower(uup.Email)
		usersEmail = strings.TrimSpace(usersEmail)
		usersEmail = strings.ToLower(usersEmail)
		// check if the email is already in the cache
		if _, ok := emailCache[usersEmail]; ok {
			log.Println("email already exists: \n", usersEmail)
			continue
		}
		if _, err := t.UserRepo.GetByEmail(usersEmail); err == nil {
			log.Printf("skipping duplicate email: %s\n", usersEmail)
			emailCache[usersEmail] = true
			continue
		}
		td, err := extractTalentData(uup)
		if err != nil {
			fmt.Println("Record", record[1])
			fmt.Printf("An error occured %s\n", err.Error())
			panic("something bad happened")
		}
		talentIndex, err := buildTalentIndex(td)
		if err != nil {
			log.Fatal(err)
		}

		manyUsers = append(manyUsers, td.User)
		manyTalents = append(manyTalents, td.Talent)
		manySkills = append(manySkills, td.Skills...)
		manyPortfolios = append(manyPortfolios, td.Portfolios...)
		manyEducations = append(manyEducations, td.Educations...)
		manyWorkExperiences = append(manyWorkExperiences, td.WorkExperiences...)

		buildSearchIndex = append(buildSearchIndex, talentIndex)

		emailCache[usersEmail] = true
	}
	var errs []error

	if len(manyUsers) < 1 {
		// TODO: this does not check if the users exists in the database,
		// it only checks if the email already exists in the emails.cache file
		// not a big deal, since we would throw away this script anyway
		fmt.Println("No users to upsert. You might want to consider deleting the cache file and resetting the database.")
		return nil
	}
	// Bulk Upsert Users
	if err := t.UserRepo.UpsertMany(manyUsers); err != nil {
		errs = append(errs, err)
		log.Println("Error upserting users:", err)
	}

	// Bulk Upsert Talents
	if err := t.TalentRepo.UpsertMany(manyTalents); err != nil {
		errs = append(errs, err)
		log.Println("Error upserting talents:", err)
	}

	// Bulk Upsert Skills
	if err := t.SkillsRepo.UpsertMany(manySkills); err != nil {
		errs = append(errs, err)
		log.Println("Error upserting skills:", err)
	}

	// Bulk Upsert Portfolios
	if err := t.PortfolioRepo.UpsertMany(manyPortfolios); err != nil {
		errs = append(errs, err)
		log.Println("Error upserting portfolios:", err)
	}

	// Bulk Upsert Educations
	if err := t.EduRepo.UpsertMany(manyEducations); err != nil {
		errs = append(errs, err)
		log.Println("Error upserting educations:", err)
	}

	// Bulk Upsert Work Experiences
	if err := t.WorkRepo.UpsertMany(manyWorkExperiences); err != nil {
		errs = append(errs, err)
		log.Println("Error upserting work experiences:", err)
	}

	// TODO: singleton pattern should be used for algolia connection
	indexStore := NewTaskMigrateSearchIndex()
	_, err = indexStore.SaveIndex(buildSearchIndex)
	if err != nil {
		return err
	}

	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}
	return nil
}

func buildTalentIndex(v *TalentData) (SearchIndexSchema, error) {
	skills := make([]Skills, 0)
	education := make([]Educations, 0)
	workExp := make([]WorkExperiences, 0)
	layoutISO := "2006-03-05"
	careerDate, err := time.Parse(layoutISO, v.Talent.ProfessionalStartDate)
	if err != nil {
		return SearchIndexSchema{}, err
	}

	talent := SearchIndexSchema{
		ID:                v.Talent.ID,
		ObjectID:          v.Talent.ID,
		Timezone:          v.Talent.TimeZone,
		IsAvailable:       v.Talent.Available,
		PreferredJobTitle: v.Talent.PreferredJobTitle,
		CareerStartDate:   careerDate,
		JobPreference:     v.Talent.JobPreference.String(),
		Country: Country{
			City:  v.Talent.City,
			Code:  v.Talent.CountryCode,
			State: v.Talent.State,
		},
	}

	for _, s := range v.Skills {
		skills = append(skills, Skills{
			Name:              s.Name,
			YearsOfExperience: int(s.YearsOfExperience),
		})
	}

	for _, s := range v.WorkExperiences {
		workExp = append(workExp, WorkExperiences{
			Location:            s.Location,
			JobTitle:            s.JobTitle,
			Description:         s.Description,
			PrimaryTechnologies: s.PrimaryTechnologies,
		})
	}

	for _, s := range v.Educations {
		education = append(education, Educations{
			Location:        s.Location,
			InstitutionName: s.InstitutionName,
			Degree:          s.Degree,
			Program:         s.Program,
		})
	}
	talent.Skills = skills
	talent.Educations = education
	talent.WorkExperiences = workExp

	return talent, nil
}

type TalentData struct {
	User            *repo.UserParams
	Talent          *repo.TalentParams
	Skills          []*repo.SkillParams
	WorkExperiences []*repo.WorkExperienceParams
	Educations      []*repo.EducationParams
	Portfolios      []*repo.PortfolioLinkParams
}

func extractTalentData(uup UpsertUserParams) (*TalentData, error) {
	user := &repo.UserParams{
		ID:        uuid.New(),
		FirstName: uup.FirstName,
		LastName:  uup.LastName,
		Email:     uup.Email,
		PhotoURL:  photo.GenerateDefaultPhoto(uup.FirstName, uup.LastName),
		Role:      userrole.Talent,
		Approved:  true,
	}

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
	var skillYears []float64
	for _, skill := range uup.Skills {
		// update skill IDs
		skill.ID = uuid.New()
		skill.TalentID = talent.ID
		skill.Preferred = true
		if skill.YearsOfExperience < 1.0 {
			skill.YearsOfExperience = 1.0
		}
		skillYears = append(skillYears, float64(skill.YearsOfExperience))
	}
	var wrkExpStartDates []string

	for _, we := range uup.WorkExperiences {
		// update work experience IDs
		we.ID = uuid.New()
		we.TalentID = talent.ID
		wrkExpStartDates = append(wrkExpStartDates, we.StartDate)
	}

	// sort an array of date strings
	sort.Strings(wrkExpStartDates)
	if len(wrkExpStartDates) > 0 {
		talent.ProfessionalStartDate = wrkExpStartDates[0]
	} else {
		sort.Float64s(skillYears)
		if len(skillYears) > 0 {
			talent.ProfessionalStartDate = time.Now().AddDate(0, 0, int(-skillYears[len(skillYears)-1])).Format("2006-01-02")
		}
	}

	date, err := time.Parse("2006-01-02", talent.ProfessionalStartDate)
	if err != nil {
		tenlog.Error(err)
		return nil, err
	}

	if date.Year() < 2000 {
		if osutil.InDevMode() {
			emailFile, err := os.OpenFile("data/reachout-emails.cache", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				tenlog.Error(err)
				return nil, err
			}
			defer emailFile.Close()
			emailFile.WriteString(user.Email + "\n")
		}
	}

	// fmt.Println("\tEducations: ")
	for _, ed := range uup.Educations {
		// update education IDs
		ed.ID = uuid.New()
		ed.TalentID = talent.ID
	}

	// fmt.Println("\tPortfolios: ")
	for _, p := range uup.Portfolios {
		// update portfolio IDs
		p.ID = uuid.New()
		p.TalentID = talent.ID
	}

	//fmt.Println("")

	fmt.Println("User ID: ", user.ID)
	fmt.Println("Talent ID: ", talent.ID)
	fmt.Print("==========================\n\n")

	return &TalentData{
		User:            user,
		Talent:          talent,
		Skills:          uup.Skills,
		Educations:      uup.Educations,
		WorkExperiences: uup.WorkExperiences,
		Portfolios:      uup.Portfolios,
	}, nil
}
