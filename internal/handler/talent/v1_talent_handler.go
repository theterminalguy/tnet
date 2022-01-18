package talent

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/internal/handler"
	repo "github.com/10hourlabs/tentn/internal/repository"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TalUUID struct {
	TalentUUID uuid.UUID `json:"talent_uuid" validate:"required"`
}

func GetCurrentTalent(c echo.Context) (*ent.Talent, error) {
	user, err := handler.GetCurrentUser(c)
	if err != nil {
		return nil, err
	}
	talent, err := repo.NewTalentRepository().GetTalentByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	return talent, nil
}

func VerifyTalentUUID(c echo.Context) (error) {
	talent, err := GetCurrentTalent(c)
	if err != nil {
		return err
	}
	req := new(TalUUID)
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &req)
	if err != nil {
		return err
	}
	// this allows read from the request body a second time
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(b))

	if req.TalentUUID != talent.UUID {
		return errors.New("unauthorized")
	}

	return nil
}
