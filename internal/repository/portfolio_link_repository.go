package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/portfoliolink"
	"github.com/google/uuid"
)

type PortfolioLinkRepository struct{}

type PortfolioLinkParams struct {
}

func NewPortfolioLinkRepository() *PortfolioLinkRepository {
	return &PortfolioLinkRepository{}
}

func (*PortfolioLinkRepository) GetAll() ([]*ent.PortfolioLink, error) {
	records, err := dBConn.PortfolioLink.Query().
		Where(portfoliolink.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*PortfolioLinkRepository) GetByUUID(id uuid.UUID) (*ent.PortfolioLink, error) {
	record, err := dBConn.PortfolioLink.Query().
		Where(portfoliolink.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, RecordNotFoundError
	}
	return record, nil
}

func (*PortfolioLinkRepository) Create(p PortfolioLinkParams) (*ent.PortfolioLink, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := dBConn.PortfolioLink.
		Create().
		// TODO: set other fields here
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *PortfolioLinkRepository) Update(id uuid.UUID, p PortfolioLinkParams) (*ent.PortfolioLink, error) {
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

func (r *PortfolioLinkRepository) DeleteByUUID(id uuid.UUID) error {
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
