package service

import (
	"errors"
	"strings"

	"github.com/theterminalguy/tentn/ent"
	repo "github.com/theterminalguy/tentn/internal/repository"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

type PortfolioLinkService struct {
	PortfolioLinkRepository *repo.PortfolioLinkRepository
	TalentRepo              *repo.TalentRepository
	UserRepo                *repo.UserRepository
	GithubService           *GithubService
}

func NewPortfolioLinkService() *PortfolioLinkService {
	return &PortfolioLinkService{
		PortfolioLinkRepository: repo.NewPortfolioLinkRepository(),
		TalentRepo:              repo.NewTalentRepository(),
		UserRepo:                repo.NewUserRepository(),
		GithubService:           NewGithubService(),
	}
}

func (h *PortfolioLinkService) Create(p repo.PortfolioLinkParams) (*ent.PortfolioLink, error) {
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

func (h *PortfolioLinkService) UpdateWithGithubProfilePicture(talentID uuid.UUID, profilePictureURL string) {
	if strings.Contains(profilePictureURL, "github.com") {
		avatarURL, err := h.GithubService.FetchUserGitHubAvatar(profilePictureURL)
		log.Print(err)
		_, vldErr := h.UserRepo.Update(talentID, repo.UserParams{PhotoURL: avatarURL})
		if vldErr != nil {
			log.Print(err)
		}
	}
}
