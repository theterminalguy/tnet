package filestorage

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/theterminalguy/tenlog"
	"github.com/theterminalguy/tentn/util"
)

type GoogleBucketFileInfo struct {
	file_path  string
	projectID  string
	bucketName string
}

const (
	BUCKET_API_URL = "https://storage.cloud.google.com"
	projectID      = "lab-internal-services"
	bucketName     = "tentn-bucket"
)

func NewGoogleBucketFileStorage(file_path string) *GoogleBucketFileInfo {
	uploader := &GoogleBucketFileInfo{
		bucketName: bucketName,
		projectID:  projectID,
		file_path:  file_path,
	}
	return uploader
}

func (c *GoogleBucketFileInfo) Upload(file io.ReadCloser) (string, error) {
	client, err := storage.NewClient(context.Background()) // TODO: avoid Background context
	if err != nil {
		return "", util.LogAndReturnErrs([]error{err}, tenlog.ERROR)
	}
	ctx := context.Background()
	uploadFolder := c.file_path
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(c.bucketName).Object(uploadFolder).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}
	link := BUCKET_API_URL + "/" + bucketName + "/" + uploadFolder

	return link, nil
}
