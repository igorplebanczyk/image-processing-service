package images

import "context"

type StorageService interface {
	UploadObject(ctx context.Context, blobName string, data []byte) error
	DownloadObject(ctx context.Context, blobName string) ([]byte, error)
	DeleteObject(ctx context.Context, blobName string) error
}
