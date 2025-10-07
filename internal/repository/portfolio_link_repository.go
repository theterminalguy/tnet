package repository

import (
	"time"

	"github.com/theterminalguy/tentn/ent"
	"github.com/theterminalguy/tentn/ent/portfoliolink"
	"github.com/theterminalguy/tentn/ent/predicate"
	"github.com/theterminalguy/tentn/util/collection"
	"github.com/google/uuid"
)

type PortfolioLinkQuerier interface {
	GetAllForTalent(talentID uuid.UUID) ([]*ent.PortfolioLink, error)
	GetAll() ([]*ent.PortfolioLink, error)
	GetByID(id uuid.UUID) (*ent.PortfolioLink, error)
	Create(p PortfolioLinkParams) (*ent.PortfolioLink, error)
	Update(id uuid.UUID, p PortfolioLinkParams) (*ent.PortfolioLink, []error)
	DeleteByID(id uuid.UUID) error
}

type PortfolioLinkRepository struct{}

type PortfolioLinkParams struct {
	ID       uuid.UUID
	URL      string    `json:"url" validate:"required,url"`
	TalentID uuid.UUID `json:"talent_id"`
	Name     string    `json:"name"`
}

func NewPortfolioLinkRepository() *PortfolioLinkRepository {
	return &PortfolioLinkRepository{}
}

func (*PortfolioLinkRepository) Filter(prd ...predicate.PortfolioLink) ([]*ent.PortfolioLink, error) {
	pfLinks, err := dBConn.PortfolioLink.Query().
		Where(prd...).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return pfLinks, nil
}

func (*PortfolioLinkRepository) GetAllForTalent(talentID uuid.UUID) ([]*ent.PortfolioLink, error) {
	records, err := dBConn.PortfolioLink.Query().
		Where(portfoliolink.And(
			portfoliolink.TalentID(talentID),
			portfoliolink.DeletedAtIsNil())).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
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

func (*PortfolioLinkRepository) GetByID(id uuid.UUID) (*ent.PortfolioLink, error) {
	record, err := dBConn.PortfolioLink.Query().
		Where(portfoliolink.ID(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*PortfolioLinkRepository) Create(p PortfolioLinkParams) (*ent.PortfolioLink, error) {
	err := ValidateParams(p)
	if err != nil {
		return nil, err
	}
	a, err := NewTalentRepository().GetByID(p.TalentID)
	if err != nil {
		return nil, err
	}
	record, err := dBConn.PortfolioLink.
		Create().
		SetTalentID(a.ID).
		SetURL(p.URL).
		SetName(p.Name).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *PortfolioLinkRepository) Update(id uuid.UUID, p PortfolioLinkParams) (*ent.PortfolioLink, []error) {
	err := ValidateParams(p, "TalentID")
	if err != nil {
		return nil, []error{err}
	}
	record, err := r.GetByID(id)
	if err != nil {
		return nil, []error{err}
	}

	var vldErrs []error
	bldr := record.Update()

	// Set and Validate URL if provided
	if vldErr := setNillableStringField(p.URL, func(v string) error {
		err := ValidateParams(p, "URL")
		if err != nil {
			return err
		}
		bldr.SetURL(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Name if provided
	if vldErr := setNillableStringField(p.Name, func(v string) error {
		err := ValidateParams(p, "Name")
		if err != nil {
			return err
		}
		bldr.SetName(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Return all validation errors at once
	// this prevents the client from making several round trips to the server
	if collection.HasAny(vldErrs) {
		return nil, vldErrs
	}

	record, err = bldr.Save(dBContext)
	if err != nil {
		return nil, []error{err}
	}

	return record, nil
}

func (r *PortfolioLinkRepository) DeleteByID(id uuid.UUID) error {
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

func (r *PortfolioLinkRepository) GetPortfolioLinkByTalentID(talentID uuid.UUID) (*ent.PortfolioLink, error) {
	record, err := dBConn.PortfolioLink.Query().
		Where(portfoliolink.And(
			portfoliolink.TalentID(talentID),
			portfoliolink.DeletedAtIsNil())).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	return record, nil
}

// UpsertMany create or update many portfolio links
func (*PortfolioLinkRepository) UpsertMany(params []*PortfolioLinkParams) error {
	builders := make([]*ent.PortfolioLinkCreate, len(params))
	for i, p := range params {
		builders[i] = dBConn.PortfolioLink.
			Create().
			SetID(p.ID).
			SetTalentID(p.TalentID).
			SetName(p.Name).
			SetURL(p.URL)
	}
	return dBConn.PortfolioLink.CreateBulk(builders...).
		Exec(dBContext)
}
