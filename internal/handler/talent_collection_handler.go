package handler

import (
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
)

type TalentCollectionHandler struct {
	TalentCollectionService    *service.TalentCollectionService
	TalentCollectionRepository *repo.TalentCollectionRepository
}

func NewTalentCollectionHandler() *TalentCollectionHandler {
	return &TalentCollectionHandler{
		TalentCollectionService:    service.NewTalentCollectionService(),
		TalentCollectionRepository: repo.NewTalentCollectionRepository(),
	}
}

func (*TalentCollectionHandler) ResourceName() string {
	return "talent_collection"
}
