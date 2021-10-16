package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/work_experience"
	"github.com/google/uuid"
)

type WorkExperienceRepository struct{}

type WorkExperienceParams struct {
}

func NewWorkExperienceRepository() *WorkExperienceRepository {
	return &WorkExperienceRepository{}
}

func (*WorkExperienceRepository) GetAll() ([]*ent.WorkExperience, error) {
	records, err := dBConn.WorkExperience.Query().
		Where(work_experience.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*WorkExperienceRepository) GetByUUID(id uuid.UUID) (*ent.WorkExperience, error) {
	record, err := dBConn.WorkExperience.Query().
		Where(work_experience.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, RecordNotFoundError
	}
	return record, nil
}

func (*WorkExperienceRepository) Create(p WorkExperienceParams) (*ent.WorkExperience, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := dBConn.WorkExperience.
		Create().
		// TODO: set other fields here
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *WorkExperienceRepository) Update(id uuid.UUID, p WorkExperienceParams) (*ent.WorkExperience, error) {
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

func (r *WorkExperienceRepository) DeleteByUUID(id uuid.UUID) error {
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
