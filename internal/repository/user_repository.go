package repository

import (
	"strings"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	"github.com/10hourlabs/tentn/ent/user"
	"github.com/10hourlabs/tentn/util"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/google/uuid"
)

type UserRepository struct{}

type UserParams struct {
	ID        uuid.UUID
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	PhotoURL  string `json:"photo_url"`
	Email     string `json:"email" validate:"required,email"`
	Role      userrole.Role
	Approved  bool
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (*UserRepository) GetAll() ([]*ent.User, error) {
	records, err := dBConn.User.Query().
		Where(user.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*UserRepository) GetByID(id uuid.UUID) (*ent.User, error) {
	record, err := dBConn.User.Query().
		Where(user.ID(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*UserRepository) GetByEmail(email string) (*ent.User, error) {
	record, err := dBConn.User.Query().
		Where(user.EmailEQ(email)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*UserRepository) Create(p UserParams) (*ent.User, error) {
	err := ValidateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := dBConn.User.
		Create().
		SetFirstName(p.FirstName).
		SetLastName(p.LastName).
		SetPhotoURL(p.PhotoURL).
		SetEmail(p.Email).
		SetRole(p.Role).
		SetApproved(p.Approved).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *UserRepository) Update(id uuid.UUID, p UserParams) (*ent.User, []error) {
	err := ValidateParams(p, "ID")
	if err != nil {
		return nil, []error{err}
	}
	record, err := r.GetByID(id)
	if err != nil {
		return nil, []error{err}
	}
	var vldErrs []error
	bldr := record.Update()

	// Set and Validate PhotoURL if provided
	if vldErr := setNillableStringField(p.PhotoURL, func(v string) error {
		err := ValidateParams(p, "PictureUrl")
		if err != nil {
			return err
		}
		bldr.SetPhotoURL(p.PhotoURL)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	if collection.HasAny(vldErrs) {
		return nil, vldErrs
	}

	record, err = bldr.Save(dBContext)
	if err != nil {
		return nil, []error{err}
	}
	return record, nil
}

func (r *UserRepository) DeleteByID(id uuid.UUID) error {
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

// UpsertMany create or update many users.
func (*UserRepository) UpsertMany(params []*UserParams) error {
	builders := make([]*ent.UserCreate, len(params))
	for i, p := range params {
		builders[i] = dBConn.User.
			Create().
			SetID(p.ID).
			SetFirstName(util.Titlelize(p.FirstName)).
			SetLastName(util.Titlelize(p.LastName)).
			SetPhotoURL(p.PhotoURL).
			SetEmail(strings.ToLower(p.Email)).
			SetRole(userrole.Talent).
			SetApproved(p.Approved)
	}
	return dBConn.User.CreateBulk(builders...).Exec(dBContext)
}

func (r *UserRepository) DeleteProfilePictureUrl(id uuid.UUID) error {
	record, err := r.GetByID(id)
	if err != nil {
		return err
	}
	_, err = record.Update().
		SetPhotoURL("").
		Save(dBContext)
	if err != nil {
		return err
	}
	return nil
}
