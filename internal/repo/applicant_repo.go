package repo

import (
	"log"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/applicant"
	"github.com/google/uuid"
)

type ApplicantRepository struct{}

type ApplicantParams struct {
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

func (p ApplicantParams) ParsedStartDate() (time.Time, error) {
	t, err := time.Parse(time.RFC3339, p.ProfessionalStartDate)
	if err != nil {
		return time.Now(), err
	}
	return t, nil
}

func NewApplicantRepository() *ApplicantRepository {
	return &ApplicantRepository{}
}

func (*ApplicantRepository) GetAll() ([]*ent.Applicant, error) {
	applicants, err := dBConn.Applicant.Query().
		Where(applicant.DeletedAtNotNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	// TODO: remove logs
	log.Println("found applicants", applicants)
	return applicants, nil
}

func (*ApplicantRepository) GetByUUID(id uuid.UUID) (*ent.Applicant, error) {
	a, err := dBConn.Applicant.Query().
		Where(applicant.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if a.DeletedAt != nil {
		return nil, RecordNotFoundError
	}
	return a, nil
}

func (*ApplicantRepository) Create(p ApplicantParams) (*ent.Applicant, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	startDate, err := p.ParsedStartDate()
	if err != nil {
		return nil, err
	}
	a, err := dBConn.Applicant.
		Create().
		SetFirstName(p.FirstName).
		SetLastName(p.LastName).
		SetPreferredName(p.PreferredName).
		SetPronoun(p.Pronoun).
		SetPreferredJobTitle(p.PreferredJobTitle).
		SetReferralCode(p.ReferralCode).
		SetProfessionalStartDate(startDate).
		SetEmail(p.Email).
		SetPhone(p.Phone).
		SetCountryCode(p.CountryCode).
		SetCity(p.City).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	// TODO: remove logs
	log.Println("applicant was created: ", a)
	return a, err
}

func (r *ApplicantRepository) Update(id uuid.UUID, p ApplicantParams) (*ent.Applicant, error) {
	a, err := r.GetByUUID(id)
	if err != nil {
		return nil, err
	}
	_, err = dBConn.Applicant.Update().
		SetEmail(p.Email).
		SetPhone(p.Phone).
		SetCountryCode(p.CountryCode).
		SetCity(p.City).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	// TODO: potential for bug
	// does a gets updated after
	// a database update? :thinking_face:
	return a, nil
}

func (r *ApplicantRepository) DeleteByUUID(id uuid.UUID) error {
	_, err := r.GetByUUID(id)
	if err != nil {
		return err
	}
	_, err = dBConn.Applicant.Update().
		SetDeletedAt(time.Now()).
		Save(dBContext)
	if err != nil {
		return err
	}
	return nil
}
