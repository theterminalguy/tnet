package repository

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/applicant"
	"github.com/10hourlabs/tentn/randutil"
	"github.com/google/uuid"
)

var invalidReferralCodeError error = errors.New("invalid referral code")

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
		Where(applicant.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return applicants, nil
}

func (*ApplicantRepository) GetByUUID(id uuid.UUID) (*ent.Applicant, error) {
	a, err := dBConn.Applicant.Query().
		Where(applicant.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	// TODO: this code server as a how to on how to
	// add edges to a node
	// ```
	// 		peeps, _ := a.QueryReferees().All(dBContext)
	// 		log.Println("Peeps", peeps)
	// 		a.Edges = ent.ApplicantEdges{
	// 			Referees: peeps,
	// 		}
	// ```
	if a.DeletedAt != nil {
		return nil, RecordNotFoundError
	}
	return a, nil
}

func (r *ApplicantRepository) Create(p ApplicantParams) (*ent.Applicant, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	startDate, err := p.ParsedStartDate()
	if err != nil {
		return nil, err
	}
	q := dBConn.Applicant.
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
		SetTentnCode(r.genTenTNCode(p)).
		SetCity(p.City)
	if len(p.ReferralCode) > 1 {
		ref, err := dBConn.Applicant.Query().
			Where(applicant.TentnCodeEQ(p.ReferralCode)).
			Only(dBContext)
		if err != nil {
			return nil, invalidReferralCodeError
		}
		q.SetReferrerID(ref.ID)
		q.SetReferralCode(ref.TentnCode)
	}
	a, err := q.Save(dBContext)
	if err != nil {
		return nil, err
	}
	return a, err
}

func (r *ApplicantRepository) Update(id uuid.UUID, p ApplicantParams) (*ent.Applicant, error) {
	err := validateParams(p, "Email", "Phone", "CountryCode", "City")
	if err != nil {
		return nil, err
	}
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
	return a, nil
}

func (r *ApplicantRepository) DeleteByUUID(id uuid.UUID) error {
	a, err := r.GetByUUID(id)
	if err != nil {
		return err
	}
	_, err = dBConn.Applicant.UpdateOne(a).
		SetDeletedAt(time.Now()).
		Save(dBContext)
	if err != nil {
		return err
	}
	return nil
}

func (r *ApplicantRepository) genTenTNCode(p ApplicantParams) string {
	attempts := 0
	for {
		if attempts == 4 {
			break
		}
		code := fmt.Sprintf("%v%v", p.PreferredName, randutil.String(5))
		a, err := dBConn.Applicant.Query().
			Where(applicant.TentnCode(code)).
			Only(dBContext)
		if a == nil {
			return code
		} else {
			log.Println(fmt.Sprintf("Duplicate TenTNCode exists: %v", err))
		}
		if err != nil {
			log.Println(err)
		}
		attempts += 1
	}
	return randutil.StringWithCharset(10, p.FirstName+p.LastName+"0123456789")
}
