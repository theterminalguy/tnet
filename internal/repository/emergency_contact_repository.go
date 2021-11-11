package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/emergencycontact"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/google/uuid"
)

type EmergencyContactRepository struct{}

type EmergencyContactParams struct {
	ApplicantUUID uuid.UUID `json:"applicant_uuid" validate:"required"`
	Name          string    `json:"name" validate:"required"`
	PhoneNumber   string    `json:"phone_number" validate:"required"`
	Address       string    `json:"address" validate:"required"`
	Relationship  string    `json:"relationship" validate:"required"`
	Email         string    `json:"email" validate:"required,email"`
}

func NewEmergencyContactRepository() *EmergencyContactRepository {
	return &EmergencyContactRepository{}
}

func (*EmergencyContactRepository) GetAll() ([]*ent.EmergencyContact, error) {
	records, err := dBConn.EmergencyContact.Query().
		Where(emergencycontact.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*EmergencyContactRepository) GetByUUID(id uuid.UUID) (*ent.EmergencyContact, error) {
	record, err := dBConn.EmergencyContact.Query().
		Where(emergencycontact.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, RecordNotFoundError
	}
	return record, nil
}

func (*EmergencyContactRepository) Create(p EmergencyContactParams) (*ent.EmergencyContact, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}

	a, err := NewApplicantRepository().GetByUUID(p.ApplicantUUID)
	if err != nil {
		return nil, err
	}

	record, err := dBConn.EmergencyContact.
		Create().
		SetApplicantID(a.ID).
		SetAddress(p.Address).
		SetEmail(p.Email).
		SetName(p.Name).
		SetPhoneNumber(p.PhoneNumber).
		SetRelationship(p.Relationship).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *EmergencyContactRepository) Update(id uuid.UUID, p EmergencyContactParams) (*ent.EmergencyContact, []error) {
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

	// Set and Validate Address if provided
	if vldErr := setNillableStringField(p.Address, func(v string) error {
		err := validateParams(p, "Address")
		if err != nil {
			return err
		}
		bldr.SetAddress(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Name if provided
	if vldErr := setNillableStringField(p.Name, func(v string) error {
		err := validateParams(p, "Name")
		if err != nil {
			return err
		}
		bldr.SetName(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate PhoneNumber if provided
	if vldErr := setNillableStringField(p.PhoneNumber, func(v string) error {
		err := validateParams(p, "PhoneNumber")
		if err != nil {
			return err
		}
		bldr.SetPhoneNumber(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Relationship if provided
	if vldErr := setNillableStringField(p.Relationship, func(v string) error {
		err := validateParams(p, "Relationship")
		if err != nil {
			return err
		}
		bldr.SetRelationship(v)
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

func (r *EmergencyContactRepository) DeleteByUUID(id uuid.UUID) error {
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
