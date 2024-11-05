package images

import "context"

type StorageService interface {
	UploadObject(ctx context.Context, objectName string, data []byte) (string, error)
	DownloadObject(ctx context.Context, url string) ([]byte, error)
	DeleteObject(ctx context.Context, url string) error
}
