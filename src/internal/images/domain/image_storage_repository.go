package domain

import "context"

type ImageStorageRepository interface {
	UploadImage(ctx context.Context, name string, bytes []byte) error
	DownloadImage(ctx context.Context, name string) ([]byte, error)
	DeleteImage(ctx context.Context, name string) error
}
