package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/schema/userrole"
	"github.com/10hourlabs/tentn/ent/user"
	"github.com/google/uuid"
)

type UserRepository struct{}

type UserParams struct {
	ID        uuid.UUID
	FirstName string        `json:"first_name" validate:"required"`
	LastName  string        `json:"last_name" validate:"required"`
	PhotoURL  string        `json:"photo_url" validate:"required"`
	Email     string        `json:"email" validate:"required,email"`
	Role      userrole.Role `json:"role"`
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
	err := validateParams(p)
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

func (r *UserRepository) Update(id uuid.UUID, p UserParams) (*ent.User, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	_, err = record.Update().
		// TODO: set other fields here
		Save(dBContext)
	if err != nil {
		return nil, err
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
func (*UserRepository) UpsertMany(params []UserParams) error {
	// TODO: should we skip existing users? i.e. users with the same email
	builders := make([]*ent.UserCreate, len(params))
	for i, p := range params {
		builders[i] = dBConn.User.
			Create().
			SetID(p.ID).
			SetFirstName(p.FirstName).
			SetLastName(p.LastName).
			SetPhotoURL(p.PhotoURL). // Fetch from Github or use default
			SetEmail(p.Email).
			SetRole(userrole.Talent).
			SetApproved(true)
	}
	//users, err := dBConn.User.CreateBulk(builders...).Save(dBContext)
	return dBConn.User.CreateBulk(builders...).
		OnConflictColumns(user.FieldEmail).
		UpdateNewValues().
		Exec(dBContext)
}
