package talent

import (
	"net/http"

	"github.com/10hourlabs/tentn/internal/service"
	"github.com/labstack/echo/v4"
)

type V1ProfilePictureHandler struct {
	PictureService service.ProfilePictureService
}

func NewV1ProfilePictureHandler() *V1ProfilePictureHandler {
	return &V1ProfilePictureHandler{
		PictureService: *service.NewProfilePictureService(),
	}
}

func (h *V1ProfilePictureHandler) Search(c echo.Context) error {
	return nil
}

func (h *V1ProfilePictureHandler) ReadAll(c echo.Context) error {
	return nil
}

func (h *V1ProfilePictureHandler) ReadByID(c echo.Context) error {
	return nil
}

func (h *V1ProfilePictureHandler) CreateOne(c echo.Context) error {
	return nil
}

func (h *V1ProfilePictureHandler) UpdateByID(c echo.Context) error {
	record, err := GetCurrentTalent(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(service.ProfilePictureParams)
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}
	params.Image = file
	vldErrs := h.PictureService.UpdateProfilePicture(record.Talent.ID, record.Talent.PhotoURL, *params)
	if vldErrs != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "updated")
}

func (h *V1ProfilePictureHandler) DeleteOne(c echo.Context) error {
	record, err := GetCurrentTalent(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.PictureService.DeleteFile(record.Talent.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
