package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/theterminalguy/tnet/ent"
	repo "github.com/theterminalguy/tnet/internal/repository"
)

const MaxNumOfFavoriteTalents = 20

var ErrFavoriteCollectionExceedsMaxNumOfTalents = fmt.Errorf(
	"favorite collection can have at most %d talents",
	MaxNumOfFavoriteTalents,
)

type TalentCollectionService struct {
	TalentCollectionRepo *repo.TalentCollectionRepository
}

func NewTalentCollectionService() *TalentCollectionService {
	return &TalentCollectionService{
		TalentCollectionRepo: repo.NewTalentCollectionRepository(),
	}
}

func (t *TalentCollectionService) AddToFavorite(id uuid.UUID, p repo.TalentCollectionParams) (*ent.TalentCollection, error) {
	collection, err := t.TalentCollectionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if len(collection.TalentUuids) >= MaxNumOfFavoriteTalents {
		return nil, ErrFavoriteCollectionExceedsMaxNumOfTalents
	}
	if len(p.TalentIDS) > MaxNumOfFavoriteTalents {
		return nil, ErrFavoriteCollectionExceedsMaxNumOfTalents
	}
	record, err := t.TalentCollectionRepo.Update(id, p)
	if err != nil {
		return nil, err
	}
	return record, nil
}
