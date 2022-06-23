package filestorage

import "mime/multipart"

type IFileStorage interface {
	Upload(file multipart.File) (string, error)
}

func NewFileStorage(storage string, file_path string) IFileStorage {

	storages := map[string]IFileStorage{
		"local":  NewLocalFileStorage(file_path),
		"google": NewGoogleBucketFileStorage(file_path),
	}

	return storages[storage]
}
