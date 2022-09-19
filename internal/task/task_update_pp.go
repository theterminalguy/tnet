package task

import (
	"errors"
	"fmt"
	"sync"

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
	var wg sync.WaitGroup
	for tId, imageUrl := range t {
		wg.Add(1)
		talentId := uuid.MustParse(tId)
		u, err := repo.NewTalentRepository().GetByID(talentId)
		if err != nil {
			tenlog.Error(fmt.Sprintf("Error getting talent %s: %s", talentId, err))
			continue
		}
		img := imageUrl
		res := fmt.Sprintf("Downloading... %s for %s", img, talentId)
		fmt.Println(res)
		go func() {
			defer wg.Done()
			link, err := service.NewTalentImageHandler(talentId, img).GetImage()
			if err != nil {
				tenlog.Error(err.Error())
				panic(err)
			}

			err = UpdateUserImage(u.UserID, link)
			if err != nil {
				tenlog.Error(err.Error())
				panic(err)
			} else {
				res := fmt.Sprintf("Downloaded... %s", link)
				fmt.Println(res)
				fmt.Println("----")
			}
		}()
	}
	wg.Wait()
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
