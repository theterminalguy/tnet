package repository

import (
	"context"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/oauth2token"
	"github.com/go-oauth2/oauth2/v4"
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

func (*Oauth2TokenRepository) GetByID(id uuid.UUID) (*ent.Oauth2Token, error) {
	record, err := dBConn.Oauth2Token.Query().
		Where(oauth2token.IDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*Oauth2TokenRepository) Create(ctx context.Context, info oauth2.TokenInfo) error {
	if ctx == nil {
		ctx = dBContext
	}
	_, err := dBConn.Oauth2Token.
		Create().
		SetUserID(uuid.MustParse(info.GetUserID())).
		SetOauth2ClientID(uuid.MustParse(info.GetClientID())).
		SetRedirectURI(info.GetRedirectURI()).
		SetScopes(info.GetScope()).
		SetCode(info.GetCode()).
		SetCodeChallenge(info.GetCodeChallenge()).
		SetCodeChallengeMethod(string(info.GetCodeChallengeMethod())).
		SetAccessToken(info.GetAccess()).
		SetRefreshToken(info.GetRefresh()).
		SetCodeCreatedAt(info.GetCodeCreateAt()).
		SetAccessTokenCreatedAt(info.GetAccessCreateAt()).
		SetRefreshTokenCreatedAt(info.GetRefreshCreateAt()).
		SetCodeExpiresIn(int64(info.GetCodeExpiresIn())).
		SetAccessTokenExpiresIn(int64(info.GetAccessExpiresIn())).
		SetRefreshTokenExpiresIn(int64(info.GetRefreshExpiresIn())).
		Save(ctx)
	return err
}

func (r *Oauth2TokenRepository) Update(id uuid.UUID, p Oauth2TokenParams) (*ent.Oauth2Token, error) {
	err := ValidateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := r.GetByID(id)
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

func (r *Oauth2TokenRepository) DeleteByID(id uuid.UUID) error {
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

// Implement the Oauth2TokenStore interface

func (r *Oauth2TokenRepository) RemoveByCode(ctx context.Context, code string) error {
	if ctx == nil {
		ctx = dBContext
	}
	record, err := dBConn.Oauth2Token.Query().
		Where(oauth2token.CodeEQ(code)).
		Only(dBContext)
	if err != nil {
		return err
	}
	if err := dBConn.Oauth2Token.DeleteOne(record).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *Oauth2TokenRepository) RemoveByAccess(ctx context.Context, access string) error {
	if ctx == nil {
		ctx = dBContext
	}
	record, err := dBConn.Oauth2Token.Query().
		Where(oauth2token.AccessTokenEQ(access)).
		Only(dBContext)
	if err != nil {
		return err
	}
	if err := dBConn.Oauth2Token.DeleteOne(record).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *Oauth2TokenRepository) RemoveByRefresh(ctx context.Context, refresh string) error {
	if ctx == nil {
		ctx = dBContext
	}
	record, err := dBConn.Oauth2Token.Query().
		Where(oauth2token.RefreshTokenEQ(refresh)).
		Only(dBContext)
	if err != nil {
		return err
	}
	if err := dBConn.Oauth2Token.DeleteOne(record).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *Oauth2TokenRepository) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	if ctx == nil {
		ctx = dBContext
	}
	record, err := dBConn.Oauth2Token.Query().
		Where(oauth2token.CodeEQ(code)).
		Only(dBContext)
	if err != nil {
		return Oauth2TokenInfo{}, err
	}
	if record.DeletedAt != nil {
		return Oauth2TokenInfo{}, ErrRecordDeleted
	}
	return Oauth2TokenInfo{
		Oauth2Token: record,
	}, nil
}

func (r *Oauth2TokenRepository) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	if ctx == nil {
		ctx = dBContext
	}
	record, err := dBConn.Oauth2Token.Query().
		Where(oauth2token.AccessTokenEQ(access)).
		Only(dBContext)
	if err != nil {
		return Oauth2TokenInfo{}, err
	}
	if record.DeletedAt != nil {
		return Oauth2TokenInfo{}, ErrRecordDeleted
	}
	return Oauth2TokenInfo{
		Oauth2Token: record,
	}, nil
}

func (r *Oauth2ClientRepository) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	if ctx == nil {
		ctx = dBContext
	}
	record, err := dBConn.Oauth2Token.Query().
		Where(oauth2token.RefreshTokenEQ(refresh)).
		Only(dBContext)
	if err != nil {
		return Oauth2TokenInfo{}, err
	}
	if record.DeletedAt != nil {
		return Oauth2TokenInfo{}, ErrRecordDeleted
	}
	return Oauth2TokenInfo{
		Oauth2Token: record,
	}, nil
}
