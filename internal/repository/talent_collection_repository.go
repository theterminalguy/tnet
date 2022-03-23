package repository

import (
	"errors"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/talentcollection"
	"github.com/google/uuid"
)

type TalentCollectionRepository struct{}

type TalentCollectionParams struct {
	UserID      int
	Name        string      `json:"name"`
	TalentUUIDS []uuid.UUID `json:"talent_uuids"`
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
	// check if a collection with the same name already exists
	records, err := dBConn.TalentCollection.Query().
		Where(
			talentcollection.And(
				talentcollection.NameEQ(p.Name),
				talentcollection.UserIDEQ(p.UserID),
			)).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	if len(records) > 0 {
		return nil, errors.New("a collection with the same name already exists")
	}

	// convert uuids to strings
	talentUUIDs := make([]string, len(p.TalentUUIDS))
	for i, uuid := range p.TalentUUIDS {
		talentUUIDs[i] = uuid.String()
	}
	record, err := dBConn.TalentCollection.
		Create().
		SetName(p.Name).
		SetUserID(p.UserID).
		SetTalentUuids(talentUUIDs).
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
	if err := dBConn.TalentCollection.DeleteOne(record).Exec(dBContext); err != nil {
		return err
	}
	return nil
}
