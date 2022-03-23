package repository

import (
	"errors"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/talentcollection"
	"github.com/10hourlabs/tentn/util/collection"
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

func (t *TalentCollectionRepository) Create(p TalentCollectionParams) (*ent.TalentCollection, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	err = t.validateScopedUniquenessOfName(p.Name, p.UserID)
	if err != nil {
		return nil, err
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
	err = r.validateScopedUniquenessOfName(p.Name, p.UserID)
	if err != nil {
		return nil, err
	}
	record, err := r.GetByUUID(id)
	if err != nil {
		return nil, err
	}
	setUUIDsForUpdate(record, p.TalentUUIDS)
	_, err = record.Update().
		SetName(p.Name).
		SetTalentUuids(record.TalentUuids).
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

func (t *TalentCollectionRepository) validateScopedUniquenessOfName(name string, userID int) error {
	records, err := dBConn.TalentCollection.Query().
		Where(
			talentcollection.And(
				talentcollection.NameEQ(name),
				talentcollection.UserIDEQ(userID),
			)).
		All(dBContext)
	if err != nil {
		return err
	}
	if len(records) > 0 {
		return errors.New("a collection with the same name already exists")
	}
	return nil
}

func setUUIDsForUpdate(t *ent.TalentCollection, newUUIDs []uuid.UUID) {
	// convert uuids to strings
	oldUUIDs := make([]uuid.UUID, len(t.TalentUuids))
	for i, talentUUID := range t.TalentUuids {
		oldUUIDs[i] = uuid.MustParse(talentUUID)
	}

	// find the uuids from the new list that are not in the old list
	newTalentUUIDs := collection.UUIDDiffs(oldUUIDs, newUUIDs)

	// if we found any new uuids, add them to the old list
	if len(newTalentUUIDs) > 0 {
		newTalentUUIDsStr := make([]string, len(newTalentUUIDs))
		for i, uuid := range newTalentUUIDs {
			newTalentUUIDsStr[i] = uuid.String()
		}
		t.TalentUuids = append(t.TalentUuids, newTalentUUIDsStr...)
	}
}
