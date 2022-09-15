package filestorage

import (
	"io"

	"github.com/10hourlabs/tenlog"
)

type FileStorager interface {
	Upload(file io.ReadCloser) (string, error)
}

func NewFileStorage(storage string, file_path string) FileStorager {
	if storage == "" || file_path == "" {
		tenlog.Error("Unable to process any storage driver")
	}
	storages := map[string]FileStorager{
		"local":  NewLocalFileStorage(file_path),
		"google": NewGoogleBucketFileStorage(file_path),
	}

	return storages[storage]
}
