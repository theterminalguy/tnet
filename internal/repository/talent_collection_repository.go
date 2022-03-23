package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/talentcollection"
	"github.com/google/uuid"
)

type TalentCollectionRepository struct{}

type TalentCollectionParams struct {
}

func NewTalentCollectionRepository() *TalentCollectionRepository {
	return &TalentCollectionRepository{}
}

func (*TalentCollectionRepository) GetAll() ([]*ent.TalentCollection, error) {
	records, err := dBConn.TalentCollection.Query().
		Where(talentcollection.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*TalentCollectionRepository) GetByUUID(id uuid.UUID) (*ent.TalentCollection, error) {
	record, err := dBConn.TalentCollection.Query().
		Where(talentcollection.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*TalentCollectionRepository) Create(p TalentCollectionParams) (*ent.TalentCollection, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := dBConn.TalentCollection.
		Create().
		// TODO: set other fields here
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *TalentCollectionRepository) Update(id uuid.UUID, p TalentCollectionParams) (*ent.TalentCollection, error) {
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

func (r *TalentCollectionRepository) DeleteByUUID(id uuid.UUID) error {
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
