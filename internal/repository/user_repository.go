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
	Name  string        `json:"name" validate:"required"`
	Email string        `json:"email" validate:"required,email"`
	Role  userrole.Role `json:"role"`
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

func (*UserRepository) GetByUUID(id uuid.UUID) (*ent.User, error) {
	record, err := dBConn.User.Query().
		Where(user.UUIDEQ(id)).
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
		SetName(p.Name).
		SetEmail(p.Email).
		SetRole(p.Role).
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
	record, err := r.GetByUUID(id)
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

func (r *UserRepository) DeleteByUUID(id uuid.UUID) error {
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
