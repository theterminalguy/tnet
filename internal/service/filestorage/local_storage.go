package filestorage

import (
	"io"
	"mime/multipart"
	"os"
)

type LocalFileInfo struct {
	file_path string
}

func NewLocalFileStorage(file_path string) *LocalFileInfo {
	return &LocalFileInfo{
		file_path: file_path,
	}
}

func (f *LocalFileInfo) Upload(file multipart.File) (string, error) {
	directory, err := os.Create(f.file_path)
	if err != nil {
		return "", err
	}
	defer directory.Close()
	if _, err := io.Copy(directory, file); err != nil {
		return "", err
	}
	return f.file_path, nil
}
