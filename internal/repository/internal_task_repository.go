package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/internaltask"
	"github.com/10hourlabs/tentn/util"
	"github.com/google/uuid"
)

type InternalTaskRepository struct{}

type InternalTaskParams struct {
	Name     string `json:"name" validate:"required"`
	Executor string `json:"executor" validate:"required,email"`
	Password string `json:"password" validate:"required"`

	// params should be a hash of key-value pairs
	Params map[string]interface{} `json:"params"`

	Succeeded bool
	Error     string
}

func NewInternalTaskRepository() *InternalTaskRepository {
	return &InternalTaskRepository{}
}

func (*InternalTaskRepository) GetAll() ([]*ent.InternalTask, error) {
	records, err := dBConn.InternalTask.Query().
		Where(internaltask.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*InternalTaskRepository) GetByID(id uuid.UUID) (*ent.InternalTask, error) {
	record, err := dBConn.InternalTask.Query().
		Where(internaltask.ID(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*InternalTaskRepository) Create(p InternalTaskParams) (*ent.InternalTask, error) {
	err := ValidateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := dBConn.InternalTask.
		Create().
		SetName(p.Name).
		SetExecutedBy(p.Executor).
		SetParams(util.MapToStringParams(p.Params)).
		SetSucceeded(p.Succeeded).
		SetError(p.Error).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *InternalTaskRepository) Update(id uuid.UUID, p InternalTaskParams) (*ent.InternalTask, error) {
	return nil, nil
}

func (r *InternalTaskRepository) DeleteByID(id uuid.UUID) error {
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
