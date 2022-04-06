package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	repo "github.com/10hourlabs/tentn/internal/repository"
	"github.com/10hourlabs/tentn/internal/repository/scope"
	"google.golang.org/api/option"
)

var (
	maxSize    int64 = 5 * 1024 * 1024 // 5MB
	projectID        = os.Getenv("PROJECT_ID")
	bucketName       = os.Getenv("BUCKET_NAME")
)

type ProfilePictureService struct {
	TalentRepo           *repo.TalentRepository
	UserRepo             *repo.UserRepository
	GoogleClientUploader *GoogleClientUploader
}

func NewProfilePictureService() *ProfilePictureService {
	return &ProfilePictureService{
		TalentRepo:           repo.NewTalentRepository(),
		UserRepo:             repo.NewUserRepository(),
		GoogleClientUploader: NewGoogleClientUploader(),
	}
}

type ProfilePictureParams struct {
	Image *multipart.FileHeader `json:"image" validate:"required"`
}

func (p *ProfilePictureService) UpdateProfilePicture(talentScope scope.TalentScope, params ProfilePictureParams) error {
	err := repo.ValidateParams(p)
	if err != nil {
		return err
	}
	err = checkFileSize(params)
	if err != nil {
		return err
	}
	err = checkFileType(params)
	if err != nil {
		return err
	}
	u, err := p.UserRepo.GetByID(talentScope.Talent.UserID)
	if err != nil {
		return err
	}
	// check if talent has a picture in google cloud storage already
	if strings.Contains(u.PhotoURL, "storage.googleapis.com") {
		// delete the old picture
		err = p.GoogleClientUploader.DeleteFile(u.PhotoURL)
		if err != nil {
			return err
		}
	}
	src, err := params.Image.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	url, err := p.GoogleClientUploader.UploadFile(src, talentScope.Talent.ID.String(), params.Image.Filename)
	if err != nil {
		return err
	}
	_, vldErr := p.UserRepo.Update(talentScope.Talent.UserID, repo.UserParams{PhotoURL: url})
	if vldErr != nil {
		return fmt.Errorf("%v", vldErr)
	}
	return nil
}

func (p *ProfilePictureService) DeleteFile(talentScope scope.TalentScope) error {
	u, err := p.UserRepo.GetByID(talentScope.Talent.UserID)
	if err != nil {
		return err
	}
	err = p.GoogleClientUploader.DeleteFile(u.PhotoURL)
	if err != nil {
		return err
	}
	err = p.UserRepo.DeleteProfilePictureUrl(talentScope.Talent.UserID)
	if err != nil {
		return err
	}
	return nil
}

type GoogleClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
}

func NewGoogleClientUploader() *GoogleClientUploader {
	client, err := storage.NewClient(context.Background(), option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))) // TODO: avoid Background context
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	uploader := &GoogleClientUploader{
		cl:         client,
		bucketName: bucketName,
		projectID:  projectID,
	}
	return uploader
}

func (c *GoogleClientUploader) UploadFile(file multipart.File, uploadPath, fileName string) (string, error) {
	ctx := context.Background()

	uploadFolder := uploadPath + "/" + fileName
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(uploadFolder).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}
	link := "https://storage.cloud.google.com" + "/" + bucketName + "/" + uploadFolder
	return link, nil
}

// deleteFile removes specified object.
func (c *GoogleClientUploader) DeleteFile(object string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	str := strings.SplitAfter(object, c.bucketName+"/")
	o := c.cl.Bucket(c.bucketName).Object(str[1])
	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", object, err)
	}
	return nil
}

func checkFileType(params ProfilePictureParams) error {
	fileType := params.Image.Header.Get("Content-Type")
	if fileType != `image/jpeg` && fileType != `image/png` || fileType == `image/jpg` {
		return fmt.Errorf("invalid file type")
	}

	return nil
}

func checkFileSize(params ProfilePictureParams) error {
	if params.Image.Size > maxSize {
		return fmt.Errorf("image size is too large")
	}
	return nil
}
