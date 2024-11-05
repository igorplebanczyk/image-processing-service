package images

import "context"

type StorageService interface {
	UploadObject(ctx context.Context, bucketName string, objectName string, data []byte) (string, error)
	GetObjectURL(bucketName string, objectName string) (string, error)
	DeleteObject(ctx context.Context, bucketName string, objectName string) error
}
