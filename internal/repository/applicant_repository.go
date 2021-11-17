package repository

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/applicant"
	"github.com/10hourlabs/tentn/randutil"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/10hourlabs/tentn/util/date"
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
	startDate, err := date.JSStringToRFC3339(p.ProfessionalStartDate)
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
		SetProfessionalStartDate(*startDate).
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

func (r *ApplicantRepository) Update(id uuid.UUID, p ApplicantParams) (*ent.Applicant, []error) {
	err := validateParams(p, "ApplicantUUID")
	if err != nil {
		return nil, []error{err}
	}
	record, err := r.GetByUUID(id)
	if err != nil {
		return nil, []error{err}
	}
	var vldErrs []error
	bldr := record.Update()

	// Set and Validate FirstName if provided
	if vldErr := setNillableStringField(p.FirstName, func(v string) error {
		err := validateParams(p, "FirsName")
		if err != nil {
			return err
		}
		bldr.SetFirstName(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate LastName if provided
	if vldErr := setNillableStringField(p.LastName, func(v string) error {
		err := validateParams(p, "LastName")
		if err != nil {
			return err
		}
		bldr.SetLastName(v)
		return nil
	}); err != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate PreferredName if provided
	if vldErr := setNillableStringField(p.PreferredName, func(program string) error {
		err := validateParams(p, "PreferredName")
		if err != nil {
			return err
		}
		bldr.SetPreferredName(p.PreferredName)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Pronoun if provided
	if vldErr := setNillableStringField(p.Pronoun, func(v string) error {
		err := validateParams(p, "Pronoun")
		if err != nil {
			return err
		}
		bldr.SetPronoun(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate PreferredJobTitle if provided
	if vldErr := setNillableStringField(p.PreferredJobTitle, func(v string) error {
		err := validateParams(p, "PreferredJobTitle")
		if err != nil {
			return err
		}
		bldr.SetPreferredJobTitle(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Email if provided
	if vldErr := setNillableStringField(p.Email, func(v string) error {
		err := validateParams(p, "Email")
		if err != nil {
			return err
		}
		bldr.SetEmail(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Phone if provided
	if vldErr := setNillableStringField(p.Phone, func(v string) error {
		err := validateParams(p, "Phone")
		if err != nil {
			return err
		}
		bldr.SetPhone(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Phone if provided
	if vldErr := setNillableStringField(p.CountryCode, func(v string) error {
		err := validateParams(p, "CountryCode")
		if err != nil {
			return err
		}
		bldr.SetCountryCode(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate City if provided
	if vldErr := setNillableStringField(p.City, func(v string) error {
		err := validateParams(p, "City")
		if err != nil {
			return err
		}
		bldr.SetCity(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Return all validation errors at once
	// this prevents the client from making several round trips to the server
	if collection.HasAny(vldErrs) {
		return nil, vldErrs
	}

	record, err = bldr.Save(dBContext)
	if err != nil {
		return nil, []error{err}
	}
	return record, nil
}

func (r *ApplicantRepository) DeleteByUUID(id uuid.UUID) error {
	record, err := r.GetByUUID(id)
	if err != nil {
		return err
	}
	_, err = record.Update().
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
