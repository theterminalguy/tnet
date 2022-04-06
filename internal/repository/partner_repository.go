package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/partner"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/google/uuid"
)

type PartnerRepository struct{}

type PartnerParams struct {
	CompanyName              string `json:"company_name" validate:"required"`
	CompanyLocation          string `json:"company_location" validate:"required"`
	WebsiteUrl               string `json:"website_url" validate:"required,url"`
	ContactPersonName        string `json:"contact_person_name" validate:"required"`
	ContactPersonPhoneNumber string `json:"contact_person_phone_number" validate:"required"`
	ContactPersonEmail       string `json:"contact_person_email" validate:"required,email"`
}

func NewPartnerRepository() *PartnerRepository {
	return &PartnerRepository{}
}

func (*PartnerRepository) GetAll() ([]*ent.Partner, error) {
	records, err := dBConn.Partner.Query().
		Where(partner.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*PartnerRepository) GetByID(id uuid.UUID) (*ent.Partner, error) {
	record, err := dBConn.Partner.Query().
		Where(partner.ID(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*PartnerRepository) Create(p PartnerParams) (*ent.Partner, error) {
	err := ValidateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := dBConn.Partner.
		Create().
		SetCompanyName(p.CompanyName).
		SetCompanyLocation(p.CompanyLocation).
		SetWebsiteUrl(p.WebsiteUrl).
		SetContactPersonPhoneNumber(p.ContactPersonPhoneNumber).
		SetContactPersonName(p.ContactPersonName).
		SetContactPersonEmail(p.ContactPersonEmail).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *PartnerRepository) Update(id uuid.UUID, p PartnerParams) (*ent.Partner, []error) {
	record, err := r.GetByID(id)
	if err != nil {
		return nil, []error{err}
	}
	var vldErrs []error
	bldr := record.Update()

	// Set and Validate CompanyName if provided
	if vldErr := setNillableStringField(p.CompanyName, func(v string) error {
		err := ValidateParams(p, "CompanyName")
		if err != nil {
			return err
		}
		bldr.SetCompanyName(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate CompanyLocation if provided
	if vldErr := setNillableStringField(p.CompanyLocation, func(v string) error {
		err := ValidateParams(p, "LastName")
		if err != nil {
			return err
		}
		bldr.SetCompanyLocation(v)
		return nil
	}); err != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate WebsiteUrl if provided
	if vldErr := setNillableStringField(p.WebsiteUrl, func(program string) error {
		err := ValidateParams(p, "WebsiteUrl")
		if err != nil {
			return err
		}
		bldr.SetWebsiteUrl(p.WebsiteUrl)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate ContactPersonName if provided
	if vldErr := setNillableStringField(p.ContactPersonName, func(v string) error {
		err := ValidateParams(p, "ContactPersonName")
		if err != nil {
			return err
		}
		bldr.SetContactPersonName(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate ContactPersonPhoneNumber if provided
	if vldErr := setNillableStringField(p.ContactPersonPhoneNumber, func(v string) error {
		err := ValidateParams(p, "ContactPersonPhoneNumber")
		if err != nil {
			return err
		}
		bldr.SetContactPersonPhoneNumber(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate ContactPersonEmail if provided
	if vldErr := setNillableStringField(p.ContactPersonEmail, func(v string) error {
		err := ValidateParams(p, "ContactPersonEmail")
		if err != nil {
			return err
		}
		bldr.SetContactPersonEmail(v)
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

func (r *PartnerRepository) DeleteByID(id uuid.UUID) error {
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
