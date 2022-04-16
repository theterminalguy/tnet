package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/oauth2token"
	"github.com/google/uuid"
)

type Oauth2TokenRepository struct{}

type Oauth2TokenParams struct {
}

func NewOauth2TokenRepository() *Oauth2TokenRepository {
	return &Oauth2TokenRepository{}
}

func (*Oauth2TokenRepository) GetAll() ([]*ent.Oauth2Token, error) {
	records, err := dBConn.Oauth2Token.Query().
		Where(oauth2token.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*Oauth2TokenRepository) GetByUUID(id uuid.UUID) (*ent.Oauth2Token, error) {
	record, err := dBConn.Oauth2Token.Query().
		Where(oauth2token.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*Oauth2TokenRepository) Create(p Oauth2TokenParams) (*ent.Oauth2Token, error) {
	err := ValidateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := dBConn.Oauth2Token.
		Create().
		// TODO: set other fields here
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *Oauth2TokenRepository) Update(id uuid.UUID, p Oauth2TokenParams) (*ent.Oauth2Token, error) {
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

func (r *Oauth2TokenRepository) DeleteByUUID(id uuid.UUID) error {
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
