package repository

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/talent"
	"github.com/10hourlabs/tentn/randutil"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/10hourlabs/tentn/util/date"
	"github.com/google/uuid"
)

var invalidReferralCodeError error = errors.New("invalid referral code")

type TalentRepository struct{}

type TalentParams struct {
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

func NewTalentRepository() *TalentRepository {
	return &TalentRepository{}
}

func (*TalentRepository) GetAll() ([]*ent.Talent, error) {
	Talents, err := dBConn.Talent.Query().
		Where(talent.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return Talents, nil
}

func (*TalentRepository) GetByUUID(id uuid.UUID) (*ent.Talent, error) {
	a, err := dBConn.Talent.Query().
		Where(talent.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	// TODO: this code server as a how to on how to
	// add edges to a node
	//```
	// peeps, _ := a.QueryReferees().All(dBContext)
	// log.Println("Peeps", peeps)
	// a.Edges = ent.ApplicantEdges{
	// 	Referees: peeps,
	// }

	pLinks, _ := a.QueryPortfoliolinks().All(dBContext)
	a.Edges = ent.TalentEdges{
		Portfoliolinks: pLinks,
	}
	//```
	if a.DeletedAt != nil {
		return nil, RecordNotFoundError
	}
	return a, nil
}

func (r *TalentRepository) Create(p TalentParams) (*ent.Talent, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	startDate, err := date.JSStringToRFC3339(p.ProfessionalStartDate)
	if err != nil {
		return nil, err
	}
	q := dBConn.Talent.
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
		ref, err := dBConn.Talent.Query().
			Where(talent.TentnCodeEQ(p.ReferralCode)).
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

func (r *TalentRepository) Update(id uuid.UUID, p TalentParams) (*ent.Talent, []error) {
	err := validateParams(p, "TalentUUID")
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

func (r *TalentRepository) DeleteByUUID(id uuid.UUID) error {
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

func (r *TalentRepository) genTenTNCode(p TalentParams) string {
	attempts := 0
	for {
		if attempts == 4 {
			break
		}
		code := fmt.Sprintf("%v%v", p.PreferredName, randutil.String(5))
		a, err := dBConn.Talent.Query().
			Where(talent.TentnCode(code)).
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
