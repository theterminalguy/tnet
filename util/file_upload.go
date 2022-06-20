package util

import (
	"fmt"
	"io"
	"os"

	"github.com/labstack/echo/v4"
)

func FileUpload(c echo.Context, filepath string, filesize int) (string, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", nil
	}
	src, err := file.Open()
	if err != nil {
		return "", nil
	}

	if file.Size > int64(filesize) {
		size := filesize / 1024 / 1024
		return fmt.Sprintf("Maximum filesize is %d%s", size, "MB"), nil
	}

	defer src.Close()
	// Move file
	file_path := fmt.Sprintf("%s/%s", filepath, file.Filename)
	directory, err := os.Create(file_path)
	if err != nil {
		return "", err
	}
	defer directory.Close()
	if _, err := io.Copy(directory, src); err != nil {
		return "", err
	}
	return file_path, nil
}
