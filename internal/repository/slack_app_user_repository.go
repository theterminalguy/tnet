package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/slackappuser"
	"github.com/google/uuid"
)

type SlackAppUserRepository struct{}

type SlackAppUserParams struct {
}

func NewSlackAppUserRepository() *SlackAppUserRepository {
	return &SlackAppUserRepository{}
}

func (*SlackAppUserRepository) GetAll() ([]*ent.SlackAppUser, error) {
	records, err := dBConn.SlackAppUser.Query().
		Where(slackappuser.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*SlackAppUserRepository) GetByUUID(id uuid.UUID) (*ent.SlackAppUser, error) {
	record, err := dBConn.SlackAppUser.Query().
		Where(slackappuser.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*SlackAppUserRepository) Create(p SlackAppUserParams) (*ent.SlackAppUser, error) {
	err := ValidateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := dBConn.SlackAppUser.
		Create().
		// TODO: set other fields here
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *SlackAppUserRepository) Update(id uuid.UUID, p SlackAppUserParams) (*ent.SlackAppUser, error) {
	err := ValidateParams(p)
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

func (r *SlackAppUserRepository) DeleteByUUID(id uuid.UUID) error {
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
