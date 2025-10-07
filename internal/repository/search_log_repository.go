package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/theterminalguy/tnet/ent"
	"github.com/theterminalguy/tnet/ent/searchlog"
)

type SearchLogRepository struct{}

type SearchLogParams struct {
	Query          string
	ResultCount    int
	Platform       string
	PlatformUserID string
	PlatformTeamID string
}

func NewSearchLogRepository() *SearchLogRepository {
	return &SearchLogRepository{}
}

func (*SearchLogRepository) GetAll() ([]*ent.SearchLog, error) {
	records, err := dBConn.SearchLog.Query().
		Where(searchlog.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*SearchLogRepository) GetByID(id uuid.UUID) (*ent.SearchLog, error) {
	record, err := dBConn.SearchLog.Query().
		Where(searchlog.ID(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*SearchLogRepository) Create(p SearchLogParams) (*ent.SearchLog, error) {
	err := ValidateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := dBConn.SearchLog.
		Create().
		SetQuery(p.Query).
		SetResultCount(p.ResultCount).
		SetPlatform(p.Platform).
		SetPlatformUserID(p.PlatformUserID).
		SetPlatformTeamID(p.PlatformTeamID).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *SearchLogRepository) Update(id uuid.UUID, p SearchLogParams) (*ent.SearchLog, error) {
	return nil, nil
}

func (r *SearchLogRepository) DeleteByID(id uuid.UUID) error {
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
