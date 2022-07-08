package recruiter

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/repository/scope"
	"github.com/10hourlabs/tentn/internal/search"
	"github.com/10hourlabs/tentn/internal/service/filestorage"
	"github.com/10hourlabs/tentn/internal/service/payment"
	"github.com/10hourlabs/tentn/oneword"
	"github.com/10hourlabs/tentn/util"
	"github.com/10hourlabs/tentn/util/osutil"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type V1RecruiterJobHandler struct {
	JobRepository           repo.JobQuerier
	JobSearch               *search.JobSearch
	JobCollectionRepository *repo.JobRepository
	FileUploadRepository    *repo.FileUploadRepository
}

func NewV1RecruiterJobHandler(jobQuerier repo.JobQuerier) *V1RecruiterJobHandler {
	return &V1RecruiterJobHandler{
		JobRepository:           jobQuerier,
		FileUploadRepository:    repo.NewFileUploadRepository(),
		JobCollectionRepository: repo.NewJobRepository(),
	}
}

func (h *V1RecruiterJobHandler) Search(c echo.Context) error {
	jobSearch := new(search.JobSearch)
	query := c.QueryString()
	records, vldErrs := jobSearch.Search(query)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (h *V1RecruiterJobHandler) ReadAll(c echo.Context) error {
	currentRecruiter := c.Get(oneword.CurrentRecruiter).(*scope.RecruiterScope)
	jobs, err := currentRecruiter.GetJobs()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, jobs)
}

// ReadByID return a job by its id. The job must be created by the recruiter
func (h *V1RecruiterJobHandler) ReadByID(c echo.Context) error {
	currentRecruiter := c.Get(oneword.CurrentRecruiter).(*scope.RecruiterScope)
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	j, err := currentRecruiter.GetJobByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, j)
}

// CreateOne creates a new job for the recruiter
func (h *V1RecruiterJobHandler) CreateOne(c echo.Context) error {
	currentRecruiter := c.Get(oneword.CurrentRecruiter).(*scope.RecruiterScope)
	recruiterID := currentRecruiter.GetID()

	const MAX_FILE_SIZE = 1024 * 1024 * 10 // 10MB
	const SUPPORTED_FILE_EXT = ".pdf"
	directory := fmt.Sprintf("job-posting/%s", recruiterID)

	file, _ := c.FormFile("file")
	src, err := file.Open()
	if err != nil {
		return err
	}

	// validate file extension
	if filepath.Ext(string(file.Filename)) != SUPPORTED_FILE_EXT {
		return c.String(http.StatusOK, fmt.Sprintf("JD only support %s File type", SUPPORTED_FILE_EXT))
	}

	// validate file size
	if file.Size > int64(MAX_FILE_SIZE) {
		size := MAX_FILE_SIZE / 1024 / 1024
		return c.String(http.StatusOK, fmt.Sprintf("Maximum filesize is %d MB", size))
	}

	// store to selected file storage
	driver := os.Getenv("FILESYSTEM_DRIVER")
	// Create directory for local storage
	if osutil.InDevMode() {
		err := os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			return err
		}
	}
	// Rename file
	now := time.Now()
	file_path := fmt.Sprintf("%s/%d%s", directory, now.UnixNano(), SUPPORTED_FILE_EXT)

	file_storage := filestorage.NewFileStorage(driver, file_path)
	path, err := file_storage.Upload(src)

	if err != nil {
		return c.String(http.StatusOK, fmt.Sprintf("Error occured %v", err))
	}
	// create DB transaction
	// TODO: Wrap the query in a Transaction (Tx)
	//create a file upload
	fileuploadParams := new(repo.FileUploadParams)
	fileuploadParams.FileUrl = path

	f, err := h.FileUploadRepository.Create(*fileuploadParams)

	if err != nil {
		return err
	}

	if f != nil {
		//create a Job Collection
		g, _ := util.SecureRandomHex(5)
		title := fmt.Sprintf("job-title-%s", g) //TODO: change this later
		params := new(repo.JobParams)
		params.Title = title
		params.Summary = "N/A"
		params.Thumbnail = "https://"
		params.WeHave = []string{""}
		params.TimeZone = "NA"
		params.Employment = "na"
		params.YouHave = []string{""}
		params.Requirements = []string{""}
		params.Category = "na"
		params.UserID = recruiterID
		params.AttachmentID = f.ID
		jd, err := h.JobRepository.Create(*params)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		}
		// Generate payment link
		driver := os.Getenv("PAYMENT_DRIVER")
		pay := payment.NewPaymentService(driver)
		_, err = pay.GenerateLink(jd.ID)
		if err != nil {
			return c.String(http.StatusBadGateway, "error occured while generate payment link")
		}
	}

	return c.String(http.StatusOK, "Document received")
}

// UpdateByID updates a job by its id. The job must be created by the recruiter
func (h *V1RecruiterJobHandler) UpdateByID(c echo.Context) error {
	currentRecruiter := c.Get(oneword.CurrentRecruiter).(*scope.RecruiterScope)
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	params := new(repo.JobParams)
	if err := c.Bind(params); err != nil {
		return err
	}
	job, err := currentRecruiter.GetJobByID(id)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	j, vldErrs := h.JobRepository.Update(job.ID, *params)
	if vldErrs != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("%v", vldErrs).Error())
	}
	return c.JSON(http.StatusOK, j)
}

// DeleteByID deletes a job by its id. The job must be created by the recruiter
func (h *V1RecruiterJobHandler) DeleteOne(c echo.Context) error {
	currentRecruiter := c.Get(oneword.CurrentRecruiter).(*scope.RecruiterScope)
	id, err := uuid.Parse(c.Param(oneword.UUID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	job, err := currentRecruiter.GetJobByID(id)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.JobRepository.DeleteByID(job.ID)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
