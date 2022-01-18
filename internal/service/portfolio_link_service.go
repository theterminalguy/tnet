package service

import (
	"errors"
	"strings"

	"github.com/10hourlabs/tentn/ent"
	repo "github.com/10hourlabs/tentn/internal/repository"
)

type PortfolioLinkService struct {
	PortfolioLinkRepository *repo.PortfolioLinkRepository
}

func NewPortfolioLinkService() *PortfolioLinkService {
	return &PortfolioLinkService{
		PortfolioLinkRepository: repo.NewPortfolioLinkRepository(),
	}
}

func (h PortfolioLinkService) Create(p repo.PortfolioLinkParams) (*ent.PortfolioLink, error) {
	// convert the name in the params to lower case
	name := strings.ToLower(p.Name)

	// check if the url contains the name
	if !strings.Contains(p.URL, name) {
		return nil, errors.New("name and url mismatch")
	}
	p.Name = name

	record, err := h.PortfolioLinkRepository.Create(p)
	if err != nil {
		return nil, err
	}

	return record, nil
}
