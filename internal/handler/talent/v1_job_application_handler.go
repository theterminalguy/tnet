package talent

import (
	"fmt"
	"net/http"

	"github.com/10hourlabs/tenlog"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/repository/scope"
	"github.com/10hourlabs/tentn/internal/search"
	"github.com/10hourlabs/tentn/internal/service"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type V1JobApplicationHandler struct {
	JobApplicationService      *service.JobApplicationService
	JobApplicationRepository   repo.JobApplicationQuerier
	TalentCollectionRepository repo.TalentCollectionRepository
}

func NewV1JobApplicationHandler(jobAppQuerier repo.JobApplicationQuerier) *V1JobApplicationHandler {
	return &V1JobApplicationHandler{
		JobApplicationService:      service.NewJobApplicationService(),
		JobApplicationRepository:   jobAppQuerier,
		TalentCollectionRepository: *repo.NewTalentCollectionRepository(),
	}
}

func (h *V1JobApplicationHandler) Search(c echo.Context) error {
	jobApplicationSearch := new(search.JobApplicationSearch)
	query := c.QueryString()
	records, vldErrs := jobApplicationSearch.Search(query)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *V1JobApplicationHandler) ReadAll(c echo.Context) error {
	currentTalent := c.Get(oneword.CurrentTalent).(*scope.TalentScope)
	records, err := currentTalent.GetJobApplications()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *V1JobApplicationHandler) ReadByID(c echo.Context) error {
	currentTalent := c.Get(oneword.CurrentTalent).(*scope.TalentScope)
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := currentTalent.GetJobApplicationByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, record)
}

func (h *V1JobApplicationHandler) CreateOne(c echo.Context) error {
	currentTalent := c.Get(oneword.CurrentTalent).(*scope.TalentScope)
	params := new(repo.JobApplicationParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params.TalentID = currentTalent.GetID()
	record, err := h.JobApplicationRepository.Create(*params)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// add talent id to job collection associated
	collection, err := h.TalentCollectionRepository.GetByJobID(params.JobUUID)
	if err != nil {
		tenlog.Error(err.Error())
		return c.String(http.StatusBadRequest, "Error occurred while processing job application")
	}
	// append new talent to collection
	talentsUuids := collection.TalentUuids
	talentsUuids = append(talentsUuids, params.TalentID.String())
	ctId := make([]uuid.UUID, 0)
	for _, tid := range talentsUuids {
		ctId = append(ctId, uuid.MustParse(tid))
	}
	collectionParams := new(repo.TalentCollectionParams)
	collectionParams.Name = collection.Name
	collectionParams.TalentIDS = ctId
	_, err = h.TalentCollectionRepository.Update(collection.ID, *collectionParams)
	if err != nil {
		tenlog.Error(err.Error())
		return c.String(http.StatusBadRequest, "Error occurred while processing job application")
	}
	return c.JSON(http.StatusCreated, record)
}

func (h *V1JobApplicationHandler) UpdateByID(c echo.Context) error {
	return nil
}

func (h *V1JobApplicationHandler) DeleteOne(c echo.Context) error {
	return nil
}
