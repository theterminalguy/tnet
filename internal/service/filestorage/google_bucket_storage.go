package filestorage

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/labstack/gommon/log"
	"google.golang.org/api/option"
)

type GoogleBucketFileInfo struct {
	file_path  string
	cl         *storage.Client
	projectID  string
	bucketName string
}

const BUCKET_API_URL = "https://storage.cloud.google.com"

var (
	projectID  = os.Getenv("PROJECT_ID")
	bucketName = os.Getenv("BUCKET_NAME")
)

func NewGoogleBucketFileStorage(file_path string) *GoogleBucketFileInfo {
	client, err := storage.NewClient(context.Background(), option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))) // TODO: avoid Background context
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	uploader := &GoogleBucketFileInfo{
		cl:         client,
		bucketName: bucketName,
		projectID:  projectID,
		file_path:  file_path,
	}
	return uploader
}

func (c *GoogleBucketFileInfo) Upload(file io.ReadCloser) (string, error) {
	ctx := context.Background()
	uploadFolder := c.file_path
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
	link := BUCKET_API_URL + "/" + bucketName + "/" + uploadFolder

	return link, nil
}
