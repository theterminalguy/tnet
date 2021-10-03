package service

import (
	"github.com/10hourlabs/tentn/ent"
	"github.com/google/uuid"
)

type CreateApplicantParams struct {
	FirstName             string `json:"first_name" validate:"required"`
	LastName              string `json:"last_name" validate:"required"`
	PreferredName         string `json:"preferred_name" validate:"required"`
	Pronoun               string `json:"pronoun" validate:"required"`
	PreferredJobTitle     string `json:"preferred_job_title" validate:"required"`
	ReferralCode          string `json:"referral_code"`
	ProfessionalStartDate string `json:"professional_start_date" validate:"datetime=2006-01-02T15:04:05Z07:00"`
	Email                 string `json:"email" validate:"required,email"`
	Phone                 string `json:"phone" validate:"required"`
	CountryCode           string `json:"country_code" validate:"required,iso3166_1_alpha2"`
	City                  string `json:"city" validate:"required"`
}

type ApplicantService struct {
}

func NewApplicantService() *ApplicantService {
	return &ApplicantService{}
}

func (*ApplicantService) Create(params *CreateApplicantParams) (*CreateApplicantParams, error) {
	if err := validateParams(params); err != nil {
		return &CreateApplicantParams{}, err
	}
	return params, nil
}
