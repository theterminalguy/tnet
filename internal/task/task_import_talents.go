package task

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/10hourlabs/tenlog"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/util/osutil"
	"github.com/10hourlabs/tentn/util/photo"
	"github.com/google/uuid"
)

type ImportTalents struct {
	UserRepo      *repo.UserRepository
	TalentRepo    *repo.TalentRepository
	SkillsRepo    *repo.SkillRepository
	EduRepo       *repo.EducationRepository
	WorkRepo      *repo.WorkExperienceRepository
	PortfolioRepo *repo.PortfolioLinkRepository
}

func NewImportTalents() *ImportTalents {
	return &ImportTalents{
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

func (t *ImportTalents) Run(_ string) error {
	manyUsers := make([]*repo.UserParams, 0)
	manyTalents := make([]*repo.TalentParams, 0)
	manySkills := make([]*repo.SkillParams, 0)
	manyPortfolios := make([]*repo.PortfolioLinkParams, 0)
	manyEducations := make([]*repo.EducationParams, 0)
	manyWorkExperiences := make([]*repo.WorkExperienceParams, 0)

	var csvReader *csv.Reader
	if os.Getenv("ENV") == "production" {
		resp, err := osutil.ReadCSVFromURL("https://storage.googleapis.com/tentn-bucket/temp/talents.csv")
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		csvReader = csv.NewReader(resp.Body)
	} else {
		f, err := os.Open("data/talents.csv")
		if err != nil {
			return err
		}
		defer f.Close()
		csvReader = csv.NewReader(f)
	}

	emailCache := make(map[string]bool)
	var emailFile *os.File
	var emailFileErr error
	if osutil.InDevMode() {
		// load a file of emails
		emailFile, emailFileErr = os.OpenFile("data/emails.cache", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if emailFileErr != nil {
			return emailFileErr
		}
		defer emailFile.Close()
		// convert the emails to a map
		scanner := bufio.NewScanner(emailFile)
		for scanner.Scan() {
			emailCache[scanner.Text()] = true
		}
	}
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		payload := []byte(record[1])
		td, err := extractTalentData(payload)
		if err != nil {
			fmt.Println("Record", record[1])
			fmt.Printf("An error occured %s\n", err.Error())
			panic("something bad happened")
		}
		talentsEmail := strings.ToLower(td.User.Email)
		if ok := emailCache[talentsEmail]; ok {
			log.Printf("skipping duplicate email: %s", talentsEmail)
			continue
		}
		// update emails map
		emailCache[talentsEmail] = true
		// write the new email to the file
		if osutil.InDevMode() {
			emailFile.WriteString(talentsEmail + "\n")
		}
		manyUsers = append(manyUsers, td.User)
		manyTalents = append(manyTalents, td.Talent)
		manySkills = append(manySkills, td.Skills...)
		manyPortfolios = append(manyPortfolios, td.Portfolios...)
		manyEducations = append(manyEducations, td.Educations...)
		manyWorkExperiences = append(manyWorkExperiences, td.WorkExperiences...)
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

	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}
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

func extractTalentData(record []byte) (*TalentData, error) {
	var uup UpsertUserParams
	err := json.Unmarshal(record, &uup)
	if err != nil {
		fmt.Println("error unmarshalling: ", err)
		return nil, err
	}
	user := &repo.UserParams{
		ID:        uuid.New(),
		FirstName: uup.FirstName,
		LastName:  uup.LastName,
		Email:     uup.Email,
		PhotoURL:  photo.GenerateDefaultPhoto(uup.FirstName, uup.LastName),
		Role:      userrole.Talent,
		Approved:  true,
	}

	/*fmt.Println(user.FirstName + " " + user.LastName)
	fmt.Println("\tID: ", user.ID)
	fmt.Println("\tFirst Name: ", user.FirstName)
	fmt.Println("\tLast Name: ", user.LastName)
	fmt.Println("\tPhoto URL: ", user.PhotoURL)
	fmt.Println("\tEmail: ", user.Email)
	fmt.Println("\tRole: ", user.Role)
	fmt.Println("\tApproved: ", user.Approved)
	fmt.Print("\t==========================\n\n")*/

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

	/*
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
	*/
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
		/*
			fmt.Println("\t\tID: ", skill.ID)
			fmt.Println("\t\tTalentID: ", skill.TalentID)
			fmt.Println("\t\tYears Of Experience: ", skill.YearsOfExperience)
			fmt.Println("\t\tPreferred: ", skill.Preferred)
			fmt.Println("\t\tNote: ", skill.Note)
			fmt.Println("\t\tName: ", skill.Name)
			fmt.Print("\t\t==========================\n\n")
		*/
	}
	var wrkExpStartDates []string

	//fmt.Println("\tWork Experiences: ")
	for _, we := range uup.WorkExperiences {
		// update work experience IDs
		we.ID = uuid.New()
		we.TalentID = talent.ID

		wrkExpStartDates = append(wrkExpStartDates, we.StartDate)
		/*
			fmt.Println("\t\tID: ", we.ID)
			fmt.Println("\t\tTalentID: ", we.TalentID)
			fmt.Println("\t\tCompany: ", we.CompanyName)
			fmt.Println("\t\tJob Title: ", we.JobTitle)
			fmt.Println("\t\tPrimary Technologies: ", we.PrimaryTechnologies)
			fmt.Println("\t\tDescription: ", we.Description)
			fmt.Println("\t\tStartDate: ", we.StartDate)
			fmt.Println("\t\tEndDate: ", we.EndDate)
			fmt.Print("\t\t==========================\n\n")
		*/
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
		/*
			fmt.Println("\t\tID: ", ed.ID)
			fmt.Println("\t\tTalentID: ", ed.TalentID)
			fmt.Println("\t\tInstitution Name: ", ed.InstitutionName)
			fmt.Println("\t\tLocation: ", ed.Location)
			fmt.Println("\t\tDegree: ", ed.Degree)
			fmt.Println("\t\tProgram: ", ed.Program)
			fmt.Println("\t\tOverview: ", ed.Overview)
			fmt.Println("\t\tStartDate: ", ed.StartDate)
			fmt.Println("\t\tEndDate: ", ed.EndDate)
			fmt.Print("\t\t==========================\n\n")
		*/
	}

	// fmt.Println("\tPortfolios: ")
	for _, p := range uup.Portfolios {
		// update portfolio IDs
		p.ID = uuid.New()
		p.TalentID = talent.ID
		/*fmt.Println("\t\tID: ", p.ID)
		fmt.Println("\t\tTalentID: ", p.TalentID)
		fmt.Println("\t\tName: ", p.Name)
		fmt.Println("\t\tURL: ", p.URL)
		fmt.Print("\t\t==========================\n\n")
		*/
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
