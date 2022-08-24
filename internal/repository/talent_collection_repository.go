package repository

import (
	"errors"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/talentcollection"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/google/uuid"
)

type TalentCollectionRepository struct{}

type TalentCollectionParams struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string      `json:"name"`
	TalentIDS []uuid.UUID `json:"talent_ids"`
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

func (*TalentCollectionRepository) GetByID(id uuid.UUID) (*ent.TalentCollection, error) {
	record, err := dBConn.TalentCollection.Query().
		Where(talentcollection.ID(id)).
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
	err := ValidateParams(p)
	if err != nil {
		return nil, err
	}
	err = t.validateScopedUniquenessOfName(p.Name, p.UserID)
	if err != nil {
		return nil, err
	}

	// convert uuids to strings
	TalentIDs := make([]string, len(p.TalentIDS))
	for i, uuid := range p.TalentIDS {
		TalentIDs[i] = uuid.String()
	}
	record, err := dBConn.TalentCollection.
		Create().
		SetName(p.Name).
		SetUserID(p.UserID).
		SetTalentUuids(TalentIDs).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *TalentCollectionRepository) Update(id uuid.UUID, p TalentCollectionParams) (*ent.TalentCollection, error) {
	err := ValidateParams(p)
	if err != nil {
		return nil, err
	}
	err = r.validateScopedUniquenessOfName(p.Name, p.UserID)
	if err != nil {
		return nil, err
	}
	record, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	prevName := record.Name
	if prevName == oneword.Favorite {
		p.Name = prevName
	}
	if p.Name == "" {
		return nil, errors.New("name is required")
	}
	setUUIDsForUpdate(record, p.TalentIDS)
	_, err = record.Update().
		SetName(p.Name).
		SetTalentUuids(record.TalentUuids).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (r *TalentCollectionRepository) DeleteByID(id uuid.UUID) error {
	record, err := r.GetByID(id)
	if err != nil {
		return err
	}
	if err := dBConn.TalentCollection.DeleteOne(record).Exec(dBContext); err != nil {
		return err
	}
	return nil
}

func (r *TalentCollectionRepository) RemoveTalents(t *ent.TalentCollection, talentIDs []uuid.UUID) (*ent.TalentCollection, error) {
	if len(t.TalentUuids) == 0 {
		return nil, errors.New("talent collection is empty")
	}
	// convert uuids to strings
	TalentIDs := make([]string, len(talentIDs))
	for i, TalentID := range talentIDs {
		TalentIDs[i] = TalentID.String()
	}
	var newTalentIDs []string
	for _, TalentID := range t.TalentUuids {
		if !collection.Contains(TalentIDs, TalentID) {
			newTalentIDs = append(newTalentIDs, TalentID)
		}
	}
	t.TalentUuids = newTalentIDs
	_, err := t.Update().
		SetTalentUuids(t.TalentUuids).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *TalentCollectionRepository) validateScopedUniquenessOfName(name string, userID uuid.UUID) error {
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
	for i, TalentID := range t.TalentUuids {
		oldUUIDs[i] = uuid.MustParse(TalentID)
	}

	// find the uuids from the new list that are not in the old list
	newTalentIDs := collection.UUIDDiffs(oldUUIDs, newUUIDs)

	// if we found any new uuids, add them to the old list
	if len(newTalentIDs) > 0 {
		newTalentIDsStr := make([]string, len(newTalentIDs))
		for i, uuid := range newTalentIDs {
			newTalentIDsStr[i] = uuid.String()
		}
		t.TalentUuids = append(t.TalentUuids, newTalentIDsStr...)
	}
}
