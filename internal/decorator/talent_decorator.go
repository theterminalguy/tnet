package decorator

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/talent"
	"github.com/google/uuid"
)

type TalentResponse struct {
	talent            *ent.Talent
	ID                uuid.UUID            `json:"id"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
	DeletedAt         *time.Time           `json:"deleted_at"`
	FirstName         string               `json:"first_name"`
	LastName          string               `json:"last_name"`
	PreferredName     string               `json:"preferred_name"`
	Pronoun           string               `json:"pronoun"`
	PreferredJobTitle string               `json:"preferred_job_title"`
	IsAvailable       bool                 `json:"is_available"`
	CareerStartDate   time.Time            `json:"career_start_date"`
	Email             string               `json:"email"`
	Phone             string               `json:"phone"`
	Country           Country              `json:"country"`
	JobPreference     talent.JobPreference `json:"job_preference"`
	Edges             *ent.TalentEdges     `json:"edges"`
	TimeZone          string               `json:"timezone"`
	ProfilePicture    string               `json:"profile_picture"`
}

func (t *TalentResponse) GetTalent() *ent.Talent {
	return t.talent
}

// TODO: decorator can accept relations to load/add to the response
func DecorateTalent(t *ent.Talent) *TalentResponse {
	return &TalentResponse{
		talent:            t,
		ID:                t.ID,
		CreatedAt:         t.CreatedAt,
		DeletedAt:         t.DeletedAt,
		FirstName:         t.FirstName,
		LastName:          t.LastName,
		PreferredName:     t.PreferredName,
		Pronoun:           t.Pronoun,
		PreferredJobTitle: t.PreferredJobTitle,
		IsAvailable:       t.IsAvailable,
		CareerStartDate:   t.ProfessionalStartDate,
		Email:             t.Email,
		Phone:             t.Phone,
		ProfilePicture:    t.Edges.User.PhotoURL,
		Country: Country{
			Code:  t.CountryCode,
			Name:  CountryLookup[t.CountryCode],
			City:  t.City,
			State: t.State,
		},
		JobPreference: t.JobPreference,
		TimeZone:      t.Timezone,
		Edges: &ent.TalentEdges{
			Educations:      t.Edges.Educations,
			WorkExperiences: t.Edges.WorkExperiences,
			Portfoliolinks:  t.Edges.Portfoliolinks,
			Skills:          t.Edges.Skills,
		},
	}
}
