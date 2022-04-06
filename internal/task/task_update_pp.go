package task

import (
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
)

type UpdateProfilePictureParams struct {
}

func NewUpdateProfilePicture() *UpdateProfilePictureParams {
	return &UpdateProfilePictureParams{}
}

func UpdateProfilePicture() error {
	allTalents, err := repo.NewPortfolioLinkRepository().GetAll()
	if err != nil {
		return err
	}

	for _, v := range allTalents {
		service.NewPortfolioLinkService().UpdateWithGithubProfilePicture(v.TalentID, v.URL)
	}
	return nil
}

func (*UpdateProfilePictureParams) Run(_ string) error {
	return UpdateProfilePicture()
}
