package task

import (
	"errors"
	"fmt"

	"github.com/10hourlabs/tenlog"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/util"
	"github.com/google/uuid"
)

type TaskUpdateProfilePicture struct {
}

func NewTaskUpdateProfilePicture() *TaskUpdateProfilePicture {
	return &TaskUpdateProfilePicture{}
}

func UpdateProfilePicture(t map[string]string) error {
	updateErrs := make([]error, 0)

	for id, img := range t {
		tid := uuid.MustParse(id)
		u, err := repo.NewTalentRepository().GetByID(tid)
		if err != nil {
			updateErrs = append(updateErrs, err)
			continue
		}
		go func(talentID uuid.UUID, imageUrl string) {
			link, err := service.NewTalentImageHandler(talentID, imageUrl).GetImage()
			if err != nil {
				updateErrs = append(updateErrs, err)
				return
			}
			err = UpdateUserImage(u.UserID, link)
			if err != nil {
				updateErrs = append(updateErrs, err)
				return
			}
			res := fmt.Sprintf("Downloaded... %s", link)
			fmt.Println(res)
			fmt.Println("----")
		}(tid, img)
	}
	if len(updateErrs) > 0 {
		return util.LogAndReturnErrs(updateErrs, tenlog.ERROR)
	}
	return nil
}

func UpdateUserImage(id uuid.UUID, f string) error {
	params := new(repo.UserParams)
	params.PhotoURL = f
	_, err := repo.NewUserRepository().Update(id, *params)
	if err != nil {
		return errors.New("error occured")
	}
	fmt.Println(f)
	return nil
}

func (*TaskUpdateProfilePicture) Run(params string) error {
	t := util.SplitStringParamsToMap(params)
	err := UpdateProfilePicture(t)
	if err != nil {
		return err
	}
	return nil
}
