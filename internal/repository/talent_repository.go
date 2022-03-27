package repository

import (
	"errors"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/predicate"
	"github.com/10hourlabs/tentn/ent/talent"
	"github.com/10hourlabs/tentn/internal/decorator"
	"github.com/10hourlabs/tentn/internal/paginator"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/10hourlabs/tentn/util/date"
	"github.com/google/uuid"
)

var ErrInvalidReferralCode error = errors.New("invalid referral code")

type TalentQuerier interface {
	GetAll() ([]*ent.Talent, error)
	GetByID(id uuid.UUID) (*ent.Talent, error)
	GetTalentByUserID(userID uuid.UUID) (*decorator.TalentResponse, error)
	Create(p TalentParams) (*decorator.TalentResponse, error)
	Update(id uuid.UUID, p TalentParams) (*decorator.TalentResponse, []error)
	DeleteByID(id uuid.UUID) error
}

type TalentRepository struct{}

type TalentParams struct {
	UserID                uuid.UUID            `json:"user_id" validate:"required"`
	FirstName             string               `json:"first_name" validate:"required"`
	LastName              string               `json:"last_name" validate:"required"`
	Email                 string               `json:"email" validate:"required"`
	PreferredName         string               `json:"preferred_name" validate:"required"`
	Pronoun               string               `json:"pronoun" validate:"required"`
	PreferredJobTitle     string               `json:"preferred_job_title" validate:"required"`
	ProfessionalStartDate string               `json:"career_start_date" validate:"required"` // YYYY-MM-DD
	Phone                 string               `json:"phone" validate:"required"`
	CountryCode           string               `json:"country_code" validate:"required,iso3166_1_alpha2"`
	City                  string               `json:"city" validate:"required"`
	JobPreference         talent.JobPreference `json:"job_preference" validate:"required"`
	Available             bool                 `json:"available"`
	TimeZone              string               `json:"timezone" validate:"required"`
	State                 string               `json:"state" validate:"required"`
}

func NewTalentRepository() *TalentRepository {
	return &TalentRepository{}
}

func (*TalentRepository) Filter(page string, prd ...predicate.Talent) (*paginator.OffsetPaginater, error) {
	// TODO: remove debug
	pager, err := paginator.NewOffsetPaginater(page)
	if err != nil {
		return nil, err
	}
	talents, err := dBConn.
		Debug().
		Talent.
		Query().
		WithUser().
		Where(prd...).
		Limit(paginator.MaxResults).
		Offset(pager.GetOffset()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	var talentList []interface{}
	for _, t := range talents {
		response := decorator.DecorateTalent(t)
		talentList = append(talentList, response)
	}
	return pager.Paginate(talentList), nil
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

func (*TalentRepository) GetByEmail(email string) (*decorator.TalentResponse, error) {
	record, err := dBConn.Talent.Query().
		Where(talent.And(
			talent.EmailEQ(email),
			talent.DeletedAtIsNil())).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	response := decorator.DecorateTalent(record)
	return response, nil
}

func (*TalentRepository) GetByID(id uuid.UUID) (*ent.Talent, error) {
	a, err := dBConn.Talent.Query().
		Where(talent.ID(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	skills, err := a.QuerySkills().All(dBContext)
	if err != nil {
		return nil, err
	}
	pLinks, _ := a.QueryPortfoliolinks().All(dBContext)
	eduLinks, _ := a.QueryEducations().All(dBContext)
	wrkExpLinks, _ := a.QueryWorkExperiences().All(dBContext)
	a.Edges = ent.TalentEdges{
		Portfoliolinks:  pLinks,
		Educations:      eduLinks,
		WorkExperiences: wrkExpLinks,
		Skills:          skills,
	}
	if a.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return a, nil
}

func (r *TalentRepository) Create(p TalentParams) (*decorator.TalentResponse, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	startDate, err := time.Parse(date.ISOLayout, p.ProfessionalStartDate)
	if err != nil {
		return nil, err
	}
	timeZoneName := TimeZoneRepo[p.TimeZone]
	if timeZoneName == nil {
		return nil, errors.New("timezone not allowed")
	}
	q := dBConn.Talent.
		Create().
		SetFirstName(p.FirstName).
		SetLastName(p.LastName).
		SetPreferredName(p.PreferredName).
		SetPronoun(p.Pronoun).
		SetPreferredJobTitle(p.PreferredJobTitle).
		SetProfessionalStartDate(startDate).
		SetEmail(p.Email).
		SetPhone(p.Phone).
		SetCountryCode(p.CountryCode).
		SetCity(p.City).
		SetUserID(p.UserID).
		SetJobPreference(p.JobPreference).
		SetIsAvailable(p.Available).
		SetTimezone(timeZoneName[1]).
		SetState(p.State)
	a, err := q.Save(dBContext)
	if err != nil {
		return nil, err
	}
	response := decorator.DecorateTalent(a)
	return response, nil
}

func (r *TalentRepository) Update(id uuid.UUID, p TalentParams) (*decorator.TalentResponse, []error) {
	err := validateParams(p, "TalentID")
	if err != nil {
		return nil, []error{err}
	}
	record, err := r.GetByID(id)
	if err != nil {
		return nil, []error{err}
	}
	var vldErrs []error
	bldr := record.Update()

	// Set and Validate FirstName if provided
	if vldErr := setNillableStringField(p.FirstName, func(v string) error {
		err := validateParams(p, "FirstName")
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

	// Set and Validate CountryCode if provided
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

	// Set and Validate JobPreference if provided
	if vldErr := setNillableStringField(string(p.JobPreference), func(v string) error {
		err := validateParams(p, "JobPreference")
		if err != nil {
			return err
		}
		bldr.SetJobPreference(p.JobPreference)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate TimeZone if provided
	if vldErr := setNillableStringField(p.TimeZone, func(v string) error {
		err := validateParams(p, "TimeZone")
		if err != nil {
			return err
		}
		timeZoneName := TimeZoneRepo[p.TimeZone]
		if timeZoneName[1] == "" {
			return errors.New("timezone not allowed")
		}
		bldr.SetTimezone(p.TimeZone)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate State if provided
	if vldErr := setNillableStringField(p.State, func(v string) error {
		err := validateParams(p, "State")
		if err != nil {
			return err
		}
		bldr.SetTimezone(p.State)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate IsAvailable if provided
	if vldErr := setNillableBoolField(p.Available, func(v bool) error {
		err := validateParams(p, "IsAvailable")
		if err != nil {
			return err
		}
		bldr.SetIsAvailable(v)
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
	response := decorator.DecorateTalent(record)
	return response, nil
}

func (r *TalentRepository) DeleteByID(id uuid.UUID) error {
	record, err := r.GetByID(id)
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

func (r *TalentRepository) GetTalentByUserID(userID uuid.UUID) (*decorator.TalentResponse, error) {
	record, err := dBConn.Talent.Query().
		WithUser().
		Where(talent.And(
			talent.UserIDEQ(userID),
			talent.DeletedAtIsNil())).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	return decorator.DecorateTalent(record), nil
}
