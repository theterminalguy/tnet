package file_upload

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/labstack/gommon/log"
	"google.golang.org/api/option"
)

type GoogleClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
}

const GOOGLE_API_URL = "https://storage.cloud.google.com"

var (
	maxSize    int64 = 5 * 1024 * 1024 // 5MB
	projectID        = os.Getenv("PROJECT_ID")
	bucketName       = os.Getenv("BUCKET_NAME")
)

type Js struct {
	Type                        string
	Project_id                  string
	Private_key_id              string
	Private_key                 string
	Client_email                string
	Client_id                   string
	Auth_uri                    string
	Token_uri                   string
	Auth_provider_x509_cert_url string
	Client_x509_cert_url        string
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
	link := GOOGLE_API_URL + "/" + bucketName + "/" + uploadFolder
	return link, nil
}
