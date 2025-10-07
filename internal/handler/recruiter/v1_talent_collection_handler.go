package recruiter

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	repo "github.com/theterminalguy/tnet/internal/repository"
	"github.com/theterminalguy/tnet/internal/repository/scope"
	"github.com/theterminalguy/tnet/internal/service"
)

type V1TalentCollectionHandler struct {
	TalentCollectionRepo    repo.TalentCollectionRepository
	TalentCollectionService service.TalentCollectionService
}

func NewV1TalentCollectionHandler() *V1TalentCollectionHandler {
	return &V1TalentCollectionHandler{
		TalentCollectionRepo:    *repo.NewTalentCollectionRepository(),
		TalentCollectionService: *service.NewTalentCollectionService(),
	}
}

func (h *V1TalentCollectionHandler) Search(c echo.Context) error {
	return nil
}

func (h *V1TalentCollectionHandler) ReadAll(c echo.Context) error {
	currentRecruiter := c.Get("currentRecruiter").(*scope.RecruiterScope)
	name := c.QueryParam("name")
	if name != "" {
		records, err := currentRecruiter.GetTalentCollectionByName(name)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, records)
	}
	records, err := currentRecruiter.GetTalentCollections()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *V1TalentCollectionHandler) ReadByID(c echo.Context) error {
	currentRecruiter := c.Get("currentRecruiter").(*scope.RecruiterScope)
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := currentRecruiter.GetTalentCollectionByID(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1TalentCollectionHandler) CreateOne(c echo.Context) error {
	currentRecruiter := c.Get("currentRecruiter").(*scope.RecruiterScope)
	params := new(repo.TalentCollectionParams)
	params.UserID = currentRecruiter.GetID()
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := h.TalentCollectionRepo.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *V1TalentCollectionHandler) UpdateByID(c echo.Context) error {
	currentRecruiter := c.Get("currentRecruiter").(*scope.RecruiterScope)
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.TalentCollectionParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	if action := c.QueryParam("action"); action == "delete" {
		// Removes talents from the collection
		record, err := currentRecruiter.DeleteTalentsFromCollection(id, params.TalentIDS)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, record)
	}

	record, err := h.TalentCollectionService.AddToFavorite(id, *params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1TalentCollectionHandler) DeleteOne(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	currentRecruiter := c.Get("currentRecruiter").(*scope.RecruiterScope)
	err = currentRecruiter.DeleteCollection(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
