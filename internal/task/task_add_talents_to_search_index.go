package task

import (
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"github.com/theterminalguy/tentn/ent"
	"github.com/theterminalguy/tentn/internal/repository"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/google/uuid"
)

type TaskMigrateSearchIndex struct {
	TalentRepo repository.TalentRepository
	initIndex  search.Index
}

type SearchIndexSchema struct {
	ID                uuid.UUID         `json:"id"`
	ObjectID          uuid.UUID         `json:"objectID"`
	PreferredJobTitle string            `json:"preferred_job_title"`
	IsAvailable       bool              `json:"is_available"`
	CareerStartDate   time.Time         `json:"career_start_date"`
	Country           Country           `json:"country"`
	JobPreference     string            `json:"job_preference"`
	Skills            []Skills          `json:"skills"`
	WorkExperiences   []WorkExperiences `json:"work_experiences"`
	Educations        []Educations      `json:"educations"`
	Timezone          string            `json:"timezone"`
}

type Country struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	City  string `json:"city"`
	State string `json:"state"`
}

type Skills struct {
	Name              string `json:"name"`
	YearsOfExperience int    `json:"years_of_experience"`
}

type WorkExperiences struct {
	Location            string   `json:"location"`
	JobTitle            string   `json:"job_title"`
	Description         string   `json:"description"`
	PrimaryTechnologies []string `json:"primary_technologies,omitempty"`
}

type Educations struct {
	InstitutionName string `json:"institution_name"`
	Location        string `json:"location"`
	Degree          string `json:"degree"`
	Program         string `json:"program"`
}

func NewTaskMigrateSearchIndex() *TaskMigrateSearchIndex {
	config := search.Configuration{
		AppID:        os.Getenv("ALGOLIA_APP_ID"),
		APIKey:       os.Getenv("ALGOLIA_WRITE_KEY"),
		MaxBatchSize: 250,
	}

	client := search.NewClientWithConfig(config)
	index := client.InitIndex("talent_search_index_schema")

	return &TaskMigrateSearchIndex{
		initIndex:  *index,
		TalentRepo: *repository.NewTalentRepository(),
	}
}

func (t *TaskMigrateSearchIndex) Run(_ string) error {
	totalTalents := t.TalentRepo.GetTotalCount()
	batchSize := 250
	batches := math.Ceil(float64(totalTalents) / float64(batchSize))

	var wg sync.WaitGroup
	for batch := 1; batch <= int(batches); batch++ {
		wg.Add(1)
		offset := (batch - 1) * batchSize
		go (func() {
			defer wg.Done()
			talents, err := t.getTalentsInBatches(offset, batchSize)
			if err != nil {
				log.Fatal(err)
			}
			_, err = t.SaveIndex(talents)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Data migrated: %v", len(talents))
			fmt.Println("")
			fmt.Println("====")
		})()
	}
	wg.Wait()

	return nil
}

func (t *TaskMigrateSearchIndex) getTalentsInBatches(offset int, limit int) ([]SearchIndexSchema, error) {
	var err error
	talents, err := t.TalentRepo.GetAllWithEdges(limit, offset)
	if err != nil {
		return nil, err
	}
	talentIndex := make([]SearchIndexSchema, 0)

	for _, t := range talents {
		talent := getSchema(t)
		talentIndex = append(talentIndex, talent)
	}
	return talentIndex, nil
}

// TODO: this should be merged with buildTalentIndex
func getSchema(v *ent.Talent) SearchIndexSchema {
	skills := make([]Skills, 0)
	education := make([]Educations, 0)
	workExp := make([]WorkExperiences, 0)
	talent := SearchIndexSchema{
		ID:                v.ID,
		ObjectID:          v.ID,
		Timezone:          v.Timezone,
		IsAvailable:       v.IsAvailable,
		PreferredJobTitle: v.PreferredJobTitle,
		CareerStartDate:   v.ProfessionalStartDate,
		JobPreference:     v.JobPreference.String(),
		Country: Country{
			City:  v.City,
			Code:  v.CountryCode,
			State: v.State,
		},
	}

	for _, s := range v.Edges.Skills {
		skills = append(skills, Skills{
			Name:              s.Name,
			YearsOfExperience: int(s.YearsOfExperience),
		})
	}

	for _, s := range v.Edges.WorkExperiences {
		workExp = append(workExp, WorkExperiences{
			Location:            s.Location,
			JobTitle:            s.JobTitle,
			Description:         s.Description,
			PrimaryTechnologies: s.PrimaryTechnologies,
		})
	}

	for _, s := range v.Edges.Educations {
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

	return talent
}

func (t *TaskMigrateSearchIndex) SaveIndex(data interface{}) (search.GroupBatchRes, error) {
	return t.initIndex.SaveObjects(data)
}
