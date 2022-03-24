package repository

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/predicate"
	"github.com/10hourlabs/tentn/ent/talent"
	"github.com/10hourlabs/tentn/internal/paginator"
	"github.com/10hourlabs/tentn/randutil"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/10hourlabs/tentn/util/date"
	"github.com/google/uuid"
)

var ErrInvalidReferralCode error = errors.New("invalid referral code")

type TalentQuerier interface {
	GetAll() ([]*ent.Talent, error)
	GetByUUID(id uuid.UUID) (*ent.Talent, error)
	GetTalentByUserID(userId int) (*ent.Talent, error)
	Create(p TalentParams) (*TalentResponse, error)
	Update(id uuid.UUID, p TalentParams) (*TalentResponse, []error)
	DeleteByUUID(id uuid.UUID) error
}

type TalentRepository struct{}

type TalentParams struct {
	UserID                int                  `json:"user_id" validate:"required"`
	FirstName             string               `json:"first_name" validate:"required"`
	LastName              string               `json:"last_name" validate:"required"`
	Email                 string               `json:"email" validate:"required"`
	PreferredName         string               `json:"preferred_name" validate:"required"`
	Pronoun               string               `json:"pronoun" validate:"required"`
	PreferredJobTitle     string               `json:"preferred_job_title" validate:"required"`
	ReferralCode          string               `json:"referral_code"`
	ProfessionalStartDate string               `json:"professional_start_date" validate:"required"` // YYYY-MM-DD
	Phone                 string               `json:"phone" validate:"required"`
	CountryCode           string               `json:"country_code" validate:"required,iso3166_1_alpha2"`
	City                  string               `json:"city" validate:"required"`
	JobPreference         talent.JobPreference `json:"job_preference" validate:"required"`
	Available             bool                 `json:"available"`
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
	talents, err := dBConn.Debug().Talent.Query().
		Where(prd...).
		Limit(paginator.MaxResults).
		Offset(pager.GetOffset()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	var talentList []interface{}
	for _, t := range talents {
		response := BuildTalentResponse(t)
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

func (*TalentRepository) GetByEmail(email string) (*TalentResponse, error) {
	record, err := dBConn.Talent.Query().
		Where(talent.And(
			talent.EmailEQ(email),
			talent.DeletedAtIsNil())).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	response := BuildTalentResponse(record)
	return response, nil
}

func (*TalentRepository) GetTalentByUUID(id uuid.UUID) (*TalentResponse, error) {
	a, err := dBConn.Talent.Query().
		Where(talent.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if a.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	skills, err := a.QuerySkills().All(dBContext)
	if err != nil {
		return nil, err
	}
	//a.Edges.Skills = append(a.Edges.Skills, skills...)
	// TODO: this code serves as a how to on how to
	// add edges to a node
	//```
	// peeps, _ := a.QueryReferees().All(dBContext)
	// log.Println("Peeps", peeps)
	// a.Edges = ent.ApplicantEdges{
	// 	Referees: peeps,
	// }

	pLinks, _ := a.QueryPortfoliolinks().All(dBContext)
	eduLinks, _ := a.QueryEducations().All(dBContext)
	wrkExpLinks, _ := a.QueryWorkExperiences().All(dBContext)

	a.Edges = ent.TalentEdges{
		Portfoliolinks:  pLinks,
		Educations:      eduLinks,
		WorkExperiences: wrkExpLinks,
		Skills:          skills,
	}
	response := BuildTalentResponse(a)

	return response, nil
}

func (*TalentRepository) GetByUUID(id uuid.UUID) (*ent.Talent, error) {
	a, err := dBConn.Talent.Query().
		Where(talent.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	skills, err := a.QuerySkills().All(dBContext)
	if err != nil {
		return nil, err
	}
	//a.Edges.Skills = append(a.Edges.Skills, skills...)
	// TODO: this code serves as a how to on how to
	// add edges to a node
	//```
	// peeps, _ := a.QueryReferees().All(dBContext)
	// log.Println("Peeps", peeps)
	// a.Edges = ent.ApplicantEdges{
	// 	Referees: peeps,
	// }

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

func (r *TalentRepository) Create(p TalentParams) (*TalentResponse, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	startDate, err := time.Parse(date.ISOLayout, p.ProfessionalStartDate)
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
		SetProfessionalStartDate(startDate).
		SetEmail(p.Email).
		SetPhone(p.Phone).
		SetCountryCode(p.CountryCode).
		SetTentnCode(r.genTenTNCode(p)).
		SetCity(p.City).
		SetUserID(p.UserID).
		SetJobPreference(p.JobPreference).
		SetIsAvailable(p.Available)
	if len(p.ReferralCode) > 1 {
		ref, err := dBConn.Talent.Query().
			Where(talent.TentnCodeEQ(p.ReferralCode)).
			Only(dBContext)
		if err != nil {
			return nil, ErrInvalidReferralCode
		}
		q.SetReferrerID(ref.ID)
		q.SetReferralCode(ref.TentnCode)
	}
	a, err := q.Save(dBContext)
	if err != nil {
		return nil, err
	}
	response := BuildTalentResponse(a)
	return response, nil
}

func (r *TalentRepository) Update(id uuid.UUID, p TalentParams) (*TalentResponse, []error) {
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

	// Set and Validate JobPreference if provided
	/*if vldErr := setNillableJSONArrayField(p.JobPreference, func(v []string) error {
		res := LinearCheckElemArray(p.JobPreference, schema.EmploymentTypes())
		if !res {
			return errors.New("unknown job preference")
		}
		err := validateParams(p, "JobPreference")
		if err != nil {
			return err
		}
		bldr.SetJobPreference(p.JobPreference)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}*/

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
	response := BuildTalentResponse(record)
	return response, nil
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

func (r *TalentRepository) GetTalentByUserID(userId int) (*ent.Talent, error) {
	record, err := dBConn.Talent.Query().
		Where(talent.And(
			talent.UserIDEQ(userId),
			talent.DeletedAtIsNil())).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	return record, nil
}

type TalentResponse struct {
	UUID                  uuid.UUID            `json:"uuid"`
	CreatedAt             time.Time            `json:"created_at"`
	UpdatedAt             time.Time            `json:"updated_at"`
	DeletedAt             *time.Time           `json:"deleted_at"`
	FirstName             string               `json:"first_name"`
	LastName              string               `json:"last_name"`
	PreferredName         string               `json:"preferred_name"`
	Pronoun               string               `json:"pronoun"`
	PreferredJobTitle     string               `json:"preferred_job_title"`
	IsAvailable           bool                 `json:"is_available"`
	ReferrerID            int                  `json:"-"`
	ReferralCode          string               `json:"referral_code"`
	TentnCode             string               `json:"tentn_code"`
	ProfessionalStartDate time.Time            `json:"professional_start_date"`
	Email                 string               `json:"email"`
	Phone                 string               `json:"phone"`
	Country               Country              `json:"country"`
	JoinedTentnAt         *time.Time           `json:"joined_tentn_at"`
	JobPreference         talent.JobPreference `json:"job_preference"`
	Edges                 *ent.TalentEdges     `json:"edges"`
}

type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
	City string `json:"city"`
}

func BuildTalentResponse(talent *ent.Talent) *TalentResponse {
	return &TalentResponse{
		UUID:                  talent.UUID,
		CreatedAt:             talent.CreatedAt,
		DeletedAt:             talent.DeletedAt,
		FirstName:             talent.FirstName,
		LastName:              talent.LastName,
		PreferredName:         talent.PreferredName,
		Pronoun:               talent.Pronoun,
		PreferredJobTitle:     talent.PreferredJobTitle,
		IsAvailable:           talent.IsAvailable,
		ReferrerID:            talent.ReferrerID,
		ReferralCode:          talent.ReferralCode,
		TentnCode:             talent.TentnCode,
		ProfessionalStartDate: talent.ProfessionalStartDate,
		Email:                 talent.Email,
		Phone:                 talent.Phone,
		Country: Country{
			Code: talent.CountryCode,
			Name: countryRepo[talent.CountryCode],
			City: talent.City,
		},
		JoinedTentnAt: talent.JoinedTentnAt,
		JobPreference: talent.JobPreference,
		Edges:         &talent.Edges,
	}
}
