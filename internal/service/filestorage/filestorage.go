package filestorage

import "mime/multipart"

type FileStorager interface {
	Upload(file multipart.File) (string, error)
}

func NewFileStorage(storage string, file_path string) FileStorager {

	storages := map[string]FileStorager{
		"local":  NewLocalFileStorage(file_path),
		"google": NewGoogleBucketFileStorage(file_path),
	}

	return storages[storage]
}
