package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

func FileUpload(c echo.Context, file_path string, filesize int) (string, error) {
	const SUPPORTED_FILE_TYPE = ".pdf"
	file, err := c.FormFile("file")
	if err != nil {
		return "", nil
	}
	src, err := file.Open()
	if err != nil {
		return "", nil
	}

	if filepath.Ext(string(file.Filename)) != SUPPORTED_FILE_TYPE {
		return fmt.Sprintf("JD only support PDF File type"), nil
	}

	if file.Size > int64(filesize) {
		size := filesize / 1024 / 1024
		return fmt.Sprintf("Maximum filesize is %d MB", size), nil
	}

	defer src.Close()
	// Move file
	now := time.Now()
	file_uploaded_path := fmt.Sprintf("%s/%d%s", file_path, now.UnixNano(), SUPPORTED_FILE_TYPE)
	directory, err := os.Create(file_uploaded_path)
	if err != nil {
		return "", err
	}
	defer directory.Close()
	if _, err := io.Copy(directory, src); err != nil {
		return "", err
	}
	return file_path, nil
}

func saveToBucket() {

}
