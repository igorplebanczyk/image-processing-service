package application

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"image"
	"image-processing-service/internal/images/application/transformations"
	"image-processing-service/internal/images/domain"
	"image/color"
	"image/png"
	"reflect"
	"testing"
	"time"
)

// Mocks

// Mocks for ImageRepository

type MockImageRepository struct {
	CreateImageFunc             func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error)
	GetImageByUserIDandNameFunc func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error)
	GetImagesByUserIDFunc       func(ctx context.Context, userID uuid.UUID, page, limit *int) ([]*domain.Image, int, error)
	UpdateImageFunc             func(ctx context.Context, id uuid.UUID) error
	DeleteImageFunc             func(ctx context.Context, id uuid.UUID) error
}

func (m *MockImageRepository) CreateImage(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
	return m.CreateImageFunc(ctx, userID, name)
}

func (m *MockImageRepository) GetImageByUserIDandName(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
	return m.GetImageByUserIDandNameFunc(ctx, userID, name)
}

func (m *MockImageRepository) GetImagesByUserID(ctx context.Context, userID uuid.UUID, page, limit *int) ([]*domain.Image, int, error) {
	return m.GetImagesByUserIDFunc(ctx, userID, page, limit)
}

func (m *MockImageRepository) UpdateImage(ctx context.Context, id uuid.UUID) error {
	return m.UpdateImageFunc(ctx, id)
}

func (m *MockImageRepository) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return m.DeleteImageFunc(ctx, id)
}

// Mocks for StorageService

type MockStorageService struct {
	UploadFunc   func(ctx context.Context, objectName string, data []byte) error
	DownloadFunc func(ctx context.Context, objectName string) ([]byte, error)
	DeleteFunc   func(ctx context.Context, objectName string) error
}

func (m *MockStorageService) Upload(ctx context.Context, objectName string, data []byte) error {
	return m.UploadFunc(ctx, objectName, data)
}

func (m *MockStorageService) Download(ctx context.Context, objectName string) ([]byte, error) {
	return m.DownloadFunc(ctx, objectName)
}

func (m *MockStorageService) Delete(ctx context.Context, objectName string) error {
	return m.DeleteFunc(ctx, objectName)
}

// Mocks for CacheService

type MockCacheService struct {
	SetFunc    func(key string, value []byte, expiration time.Duration) error
	GetFunc    func(key string) ([]byte, error)
	DeleteFunc func(key string) error
}

func (m *MockCacheService) Set(key string, value []byte, expiration time.Duration) error {
	return m.SetFunc(key, value, expiration)
}

func (m *MockCacheService) Get(key string) ([]byte, error) {
	return m.GetFunc(key)
}

func (m *MockCacheService) Delete(key string) error {
	return m.DeleteFunc(key)
}

// Tests

func TestImageService_UploadImage(t *testing.T) {
	validUUID := uuid.New()
	validImageName := "valid_image.jpg"
	validImageBytes := []byte{1, 2, 3}

	type fields struct {
		repo                  domain.ImageRepository
		storage               domain.StorageService
		cache                 domain.CacheService
		transformationService *transformations.Service
	}
	type args struct {
		userID     uuid.UUID
		imageName  string
		imageBytes []byte
	}

	expectedImage := &domain.Image{
		ID:     uuid.New(),
		UserID: validUUID,
		Name:   validImageName,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Image
		wantErr bool
	}{
		{
			name: "Successful upload",
			fields: fields{
				repo: &MockImageRepository{
					CreateImageFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return expectedImage, nil
					},
				},
				storage: &MockStorageService{
					UploadFunc: func(ctx context.Context, objectName string, data []byte) error {
						return nil
					},
				},
				cache: &MockCacheService{
					SetFunc: func(key string, value []byte, expiration time.Duration) error {
						return nil
					},
				},
			},
			args: args{
				userID:     validUUID,
				imageName:  validImageName,
				imageBytes: validImageBytes,
			},
			want:    expectedImage,
			wantErr: false,
		},
		{
			name: "Invalid image name",
			fields: fields{
				repo: &MockImageRepository{
					CreateImageFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return nil, fmt.Errorf("should not be called")
					},
				},
				storage: &MockStorageService{
					UploadFunc: func(ctx context.Context, objectName string, data []byte) error {
						return fmt.Errorf("should not be called")
					},
				},
				cache: &MockCacheService{
					SetFunc: func(key string, value []byte, expiration time.Duration) error {
						return fmt.Errorf("should not be called")
					},
				},
			},
			args: args{
				userID:     validUUID,
				imageName:  "",
				imageBytes: validImageBytes,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid image data",
			fields: fields{
				repo: &MockImageRepository{
					CreateImageFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return nil, fmt.Errorf("should not be called")
					},
				},
				storage: &MockStorageService{
					UploadFunc: func(ctx context.Context, objectName string, data []byte) error {
						return fmt.Errorf("should not be called")
					},
				},
				cache: &MockCacheService{
					SetFunc: func(key string, value []byte, expiration time.Duration) error {
						return fmt.Errorf("should not be called")
					},
				},
			},
			args: args{
				userID:     validUUID,
				imageName:  validImageName,
				imageBytes: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Repository error on CreateImage",
			fields: fields{
				repo: &MockImageRepository{
					CreateImageFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return nil, fmt.Errorf("repository error")
					},
				},
				storage: &MockStorageService{
					UploadFunc: func(ctx context.Context, objectName string, data []byte) error {
						return nil
					},
				},
				cache: &MockCacheService{
					SetFunc: func(key string, value []byte, expiration time.Duration) error {
						return nil
					},
				},
			},
			args: args{
				userID:     validUUID,
				imageName:  validImageName,
				imageBytes: validImageBytes,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Storage upload error",
			fields: fields{
				repo: &MockImageRepository{
					CreateImageFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return expectedImage, nil
					},
				},
				storage: &MockStorageService{
					UploadFunc: func(ctx context.Context, objectName string, data []byte) error {
						return fmt.Errorf("storage error")
					},
				},
				cache: &MockCacheService{
					SetFunc: func(key string, value []byte, expiration time.Duration) error {
						return nil
					},
				},
			},
			args: args{
				userID:     validUUID,
				imageName:  validImageName,
				imageBytes: validImageBytes,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Cache set error",
			fields: fields{
				repo: &MockImageRepository{
					CreateImageFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return expectedImage, nil
					},
				},
				storage: &MockStorageService{
					UploadFunc: func(ctx context.Context, objectName string, data []byte) error {
						return nil
					},
				},
				cache: &MockCacheService{
					SetFunc: func(key string, value []byte, expiration time.Duration) error {
						return fmt.Errorf("cache error")
					},
				},
			},
			args: args{
				userID:     validUUID,
				imageName:  validImageName,
				imageBytes: validImageBytes,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ImageService{
				repo:                  tt.fields.repo,
				storage:               tt.fields.storage,
				cache:                 tt.fields.cache,
				transformationService: tt.fields.transformationService,
			}
			got, err := s.UploadImage(tt.args.userID, tt.args.imageName, tt.args.imageBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UploadImage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageService_ListUserImages(t *testing.T) {
	validUUID := uuid.New()
	page := 1
	limit := 10

	type fields struct {
		repo                  domain.ImageRepository
		storage               domain.StorageService
		cache                 domain.CacheService
		transformationService *transformations.Service
	}
	type args struct {
		userID uuid.UUID
		page   *int
		limit  *int
	}

	mockImages := []*domain.Image{
		{ID: uuid.New(), UserID: validUUID, Name: "image1.jpg"},
		{ID: uuid.New(), UserID: validUUID, Name: "image2.jpg"},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*domain.Image
		want1   int
		wantErr bool
	}{
		{
			name: "Successful list of images",
			fields: fields{
				repo: &MockImageRepository{
					GetImagesByUserIDFunc: func(ctx context.Context, userID uuid.UUID, page, limit *int) ([]*domain.Image, int, error) {
						return mockImages, 2, nil
					},
				},
				storage: &MockStorageService{},
				cache:   &MockCacheService{},
			},
			args: args{
				userID: validUUID,
				page:   &page,
				limit:  &limit,
			},
			want:    mockImages,
			want1:   2,
			wantErr: false,
		},
		{
			name: "Repository error",
			fields: fields{
				repo: &MockImageRepository{
					GetImagesByUserIDFunc: func(ctx context.Context, userID uuid.UUID, page, limit *int) ([]*domain.Image, int, error) {
						return nil, 0, fmt.Errorf("repository error")
					},
				},
				storage: &MockStorageService{},
				cache:   &MockCacheService{},
			},
			args: args{
				userID: validUUID,
				page:   &page,
				limit:  &limit,
			},
			want:    nil,
			want1:   0,
			wantErr: true,
		},
		{
			name: "Empty image list with valid pagination",
			fields: fields{
				repo: &MockImageRepository{
					GetImagesByUserIDFunc: func(ctx context.Context, userID uuid.UUID, page, limit *int) ([]*domain.Image, int, error) {
						return []*domain.Image{}, 0, nil
					},
				},
				storage: &MockStorageService{},
				cache:   &MockCacheService{},
			},
			args: args{
				userID: validUUID,
				page:   &page,
				limit:  &limit,
			},
			want:    []*domain.Image{},
			want1:   0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ImageService{
				repo:                  tt.fields.repo,
				storage:               tt.fields.storage,
				cache:                 tt.fields.cache,
				transformationService: tt.fields.transformationService,
			}
			got, got1, err := s.ListUserImages(tt.args.userID, tt.args.page, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListUserImages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListUserImages() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ListUserImages() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestImageService_GetImageData(t *testing.T) {
	validUUID := uuid.New()
	validImageName := "test_image.jpg"

	type fields struct {
		repo                  domain.ImageRepository
		storage               domain.StorageService
		cache                 domain.CacheService
		transformationService *transformations.Service
	}
	type args struct {
		userID    uuid.UUID
		imageName string
	}

	mockImage := &domain.Image{
		ID:     uuid.New(),
		UserID: validUUID,
		Name:   validImageName,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Image
		wantErr bool
	}{
		{
			name: "Successful image retrieval",
			fields: fields{
				repo: &MockImageRepository{
					GetImageByUserIDandNameFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return mockImage, nil
					},
				},
				storage: &MockStorageService{},
				cache:   &MockCacheService{},
			},
			args: args{
				userID:    validUUID,
				imageName: validImageName,
			},
			want:    mockImage,
			wantErr: false,
		},
		{
			name: "Image not found",
			fields: fields{
				repo: &MockImageRepository{
					GetImageByUserIDandNameFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return nil, fmt.Errorf("image not found")
					},
				},
				storage: &MockStorageService{},
				cache:   &MockCacheService{},
			},
			args: args{
				userID:    validUUID,
				imageName: "nonexistent_image.jpg",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Repository error",
			fields: fields{
				repo: &MockImageRepository{
					GetImageByUserIDandNameFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return nil, fmt.Errorf("repository error")
					},
				},
				storage: &MockStorageService{},
				cache:   &MockCacheService{},
			},
			args: args{
				userID:    validUUID,
				imageName: validImageName,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ImageService{
				repo:                  tt.fields.repo,
				storage:               tt.fields.storage,
				cache:                 tt.fields.cache,
				transformationService: tt.fields.transformationService,
			}
			got, err := s.GetImageData(tt.args.userID, tt.args.imageName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetImageData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetImageData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageService_DownloadImage(t *testing.T) {
	validUUID := uuid.New()
	validImageName := "test_image.jpg"
	mockImageData := []byte{1, 2, 3, 4, 5}

	type fields struct {
		repo                  domain.ImageRepository
		storage               domain.StorageService
		cache                 domain.CacheService
		transformationService *transformations.Service
	}

	type args struct {
		userID    uuid.UUID
		imageName string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Successful retrieval from cache",
			fields: fields{
				cache: &MockCacheService{
					GetFunc: func(key string) ([]byte, error) {
						return mockImageData, nil
					},
				},
				storage: &MockStorageService{
					DownloadFunc: func(ctx context.Context, objectName string) ([]byte, error) {
						return nil, fmt.Errorf("should not reach storage if cache hit")
					},
				},
			},
			args: args{
				userID:    validUUID,
				imageName: validImageName,
			},
			want:    mockImageData,
			wantErr: false,
		},
		{
			name: "Cache miss, successful retrieval from storage",
			fields: fields{
				cache: &MockCacheService{
					GetFunc: func(key string) ([]byte, error) {
						return nil, fmt.Errorf("cache miss")
					},
				},
				storage: &MockStorageService{
					DownloadFunc: func(ctx context.Context, objectName string) ([]byte, error) {
						return mockImageData, nil
					},
				},
			},
			args: args{
				userID:    validUUID,
				imageName: validImageName,
			},
			want:    mockImageData,
			wantErr: false,
		},
		{
			name: "Cache miss and storage error",
			fields: fields{
				cache: &MockCacheService{
					GetFunc: func(key string) ([]byte, error) {
						return nil, fmt.Errorf("cache miss")
					},
				},
				storage: &MockStorageService{
					DownloadFunc: func(ctx context.Context, objectName string) ([]byte, error) {
						return nil, fmt.Errorf("storage error")
					},
				},
			},
			args: args{
				userID:    validUUID,
				imageName: validImageName,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ImageService{
				repo:                  tt.fields.repo,
				storage:               tt.fields.storage,
				cache:                 tt.fields.cache,
				transformationService: tt.fields.transformationService,
			}
			got, err := s.DownloadImage(tt.args.userID, tt.args.imageName)
			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DownloadImage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageService_DeleteImage(t *testing.T) {
	validUUID := uuid.New()
	validImageName := "test_image.jpg"
	imageID := uuid.New()

	type fields struct {
		repo                  domain.ImageRepository
		storage               domain.StorageService
		cache                 domain.CacheService
		transformationService *transformations.Service
	}

	type args struct {
		userID    uuid.UUID
		imageName string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Successful deletion",
			fields: fields{
				repo: &MockImageRepository{
					GetImageByUserIDandNameFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return &domain.Image{ID: imageID}, nil
					},
					DeleteImageFunc: func(ctx context.Context, id uuid.UUID) error {
						return nil
					},
				},
				storage: &MockStorageService{
					DeleteFunc: func(ctx context.Context, objectName string) error {
						return nil
					},
				},
				cache: &MockCacheService{
					DeleteFunc: func(key string) error {
						return nil
					},
				},
			},
			args: args{
				userID:    validUUID,
				imageName: validImageName,
			},
			wantErr: false,
		},
		{
			name: "Error retrieving image ID from repo",
			fields: fields{
				repo: &MockImageRepository{
					GetImageByUserIDandNameFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return nil, fmt.Errorf("error fetching image")
					},
					DeleteImageFunc: func(ctx context.Context, id uuid.UUID) error {
						return nil
					},
				},
				storage: &MockStorageService{
					DeleteFunc: func(ctx context.Context, objectName string) error {
						return nil
					},
				},
				cache: &MockCacheService{
					DeleteFunc: func(key string) error {
						return nil
					},
				},
			},
			args: args{
				userID:    validUUID,
				imageName: validImageName,
			},
			wantErr: true,
		},
		{
			name: "Error deleting image from repo",
			fields: fields{
				repo: &MockImageRepository{
					GetImageByUserIDandNameFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return &domain.Image{ID: imageID}, nil
					},
					DeleteImageFunc: func(ctx context.Context, id uuid.UUID) error {
						return fmt.Errorf("error deleting from repo")
					},
				},
				storage: &MockStorageService{
					DeleteFunc: func(ctx context.Context, objectName string) error {
						return nil
					},
				},
				cache: &MockCacheService{
					DeleteFunc: func(key string) error {
						return nil
					},
				},
			},
			args: args{
				userID:    validUUID,
				imageName: validImageName,
			},
			wantErr: true,
		},
		{
			name: "Error deleting image from storage",
			fields: fields{
				repo: &MockImageRepository{
					GetImageByUserIDandNameFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return &domain.Image{ID: imageID}, nil
					},
					DeleteImageFunc: func(ctx context.Context, id uuid.UUID) error {
						return nil
					},
				},
				storage: &MockStorageService{
					DeleteFunc: func(ctx context.Context, objectName string) error {
						return fmt.Errorf("error deleting from storage")
					},
				},
				cache: &MockCacheService{
					DeleteFunc: func(key string) error {
						return nil
					},
				},
			},
			args: args{
				userID:    validUUID,
				imageName: validImageName,
			},
			wantErr: true,
		},
		{
			name: "Error deleting image from cache",
			fields: fields{
				repo: &MockImageRepository{
					GetImageByUserIDandNameFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return &domain.Image{ID: imageID}, nil
					},
					DeleteImageFunc: func(ctx context.Context, id uuid.UUID) error {
						return nil
					},
				},
				storage: &MockStorageService{
					DeleteFunc: func(ctx context.Context, objectName string) error {
						return nil
					},
				},
				cache: &MockCacheService{
					DeleteFunc: func(key string) error {
						return fmt.Errorf("error deleting from cache")
					},
				},
			},
			args: args{
				userID:    validUUID,
				imageName: validImageName,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ImageService{
				repo:                  tt.fields.repo,
				storage:               tt.fields.storage,
				cache:                 tt.fields.cache,
				transformationService: tt.fields.transformationService,
			}
			if err := s.DeleteImage(tt.args.userID, tt.args.imageName); (err != nil) != tt.wantErr {
				t.Errorf("DeleteImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImageService_ApplyTransformations(t *testing.T) {
	validUUID := uuid.New()
	validImageName := "test_image.png"
	mockImageBytes := generatePNGImage()

	imageID := uuid.New()
	transformation1 := domain.Transformation{
		Type: domain.Resize,
		Options: map[string]float64{
			"width":  100,
			"height": 100,
		},
	}
	transformation2 := domain.Transformation{
		Type:    domain.Grayscale,
		Options: map[string]float64{},
	}
	transformationsList := []domain.Transformation{transformation1, transformation2}

	type fields struct {
		repo                  domain.ImageRepository
		storage               domain.StorageService
		cache                 domain.CacheService
		transformationService *transformations.Service
	}

	type args struct {
		userID          uuid.UUID
		imageName       string
		transformations []domain.Transformation
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Successful transformations and upload",
			fields: fields{
				repo: &MockImageRepository{
					GetImageByUserIDandNameFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return &domain.Image{ID: imageID}, nil
					},
					UpdateImageFunc: func(ctx context.Context, id uuid.UUID) error {
						return nil
					},
				},
				storage: &MockStorageService{
					DownloadFunc: func(ctx context.Context, objectName string) ([]byte, error) {
						return mockImageBytes, nil
					},
					UploadFunc: func(ctx context.Context, objectName string, data []byte) error {
						return nil
					},
				},
				cache: &MockCacheService{
					GetFunc: func(key string) ([]byte, error) {
						return nil, fmt.Errorf("cache miss")
					},
					SetFunc: func(key string, data []byte, duration time.Duration) error {
						return nil
					},
				},
				transformationService: transformations.New(10, 100),
			},
			args: args{
				userID:          validUUID,
				imageName:       validImageName,
				transformations: transformationsList,
			},
			wantErr: false,
		},
		{
			name: "Download image error",
			fields: fields{
				repo: &MockImageRepository{},
				storage: &MockStorageService{
					DownloadFunc: func(ctx context.Context, objectName string) ([]byte, error) {
						return nil, fmt.Errorf("error downloading image")
					},
					UploadFunc: func(ctx context.Context, objectName string, data []byte) error {
						return nil
					},
				},
				cache: &MockCacheService{
					GetFunc: func(key string) ([]byte, error) {
						return nil, fmt.Errorf("cache miss")
					},
					SetFunc: func(key string, data []byte, duration time.Duration) error {
						return nil
					},
				},
				transformationService: transformations.New(10, 100),
			},
			args: args{
				userID:          validUUID,
				imageName:       "nonexistent_image.jpg",
				transformations: transformationsList,
			},
			wantErr: true,
		},
		{
			name: "Successful transformation with cache hit",
			fields: fields{
				repo: &MockImageRepository{
					GetImageByUserIDandNameFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return &domain.Image{ID: imageID}, nil
					},
					UpdateImageFunc: func(ctx context.Context, id uuid.UUID) error {
						return nil
					},
				},
				storage: &MockStorageService{
					UploadFunc: func(ctx context.Context, objectName string, data []byte) error {
						return nil
					},
				},
				cache: &MockCacheService{
					GetFunc: func(key string) ([]byte, error) {
						return mockImageBytes, nil
					},
					SetFunc: func(key string, data []byte, duration time.Duration) error {
						return nil
					},
				},
				transformationService: transformations.New(10, 100),
			},
			args: args{
				userID:          validUUID,
				imageName:       validImageName,
				transformations: transformationsList,
			},
			wantErr: false,
		},
		{
			name: "Error caching transformed image",
			fields: fields{
				repo: &MockImageRepository{
					GetImageByUserIDandNameFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return &domain.Image{ID: imageID}, nil
					},
					UpdateImageFunc: func(ctx context.Context, id uuid.UUID) error {
						return nil
					},
				},
				storage: &MockStorageService{
					DownloadFunc: func(ctx context.Context, objectName string) ([]byte, error) {
						return mockImageBytes, nil
					},
					UploadFunc: func(ctx context.Context, objectName string, data []byte) error {
						return nil
					},
				},
				cache: &MockCacheService{
					GetFunc: func(key string) ([]byte, error) {
						return nil, fmt.Errorf("cache miss")
					},
					SetFunc: func(key string, data []byte, duration time.Duration) error {
						return fmt.Errorf("error caching image")
					},
				},
				transformationService: transformations.New(10, 100),
			},
			args: args{
				userID:          validUUID,
				imageName:       validImageName,
				transformations: transformationsList,
			},
			wantErr: true,
		},
		{
			name: "Error updating image metadata in repository",
			fields: fields{
				repo: &MockImageRepository{
					GetImageByUserIDandNameFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return &domain.Image{ID: imageID}, nil
					},
					UpdateImageFunc: func(ctx context.Context, id uuid.UUID) error {
						return fmt.Errorf("error updating image")
					},
				},
				storage: &MockStorageService{
					DownloadFunc: func(ctx context.Context, objectName string) ([]byte, error) {
						return mockImageBytes, nil
					},
					UploadFunc: func(ctx context.Context, objectName string, data []byte) error {
						return nil
					},
				},
				cache: &MockCacheService{
					GetFunc: func(key string) ([]byte, error) {
						return nil, nil
					},
					SetFunc: func(key string, data []byte, duration time.Duration) error {
						return nil
					},
				},
				transformationService: transformations.New(10, 100),
			},
			args: args{
				userID:          validUUID,
				imageName:       validImageName,
				transformations: transformationsList,
			},
			wantErr: true,
		},
		{
			name: "Error uploading transformed image",
			fields: fields{
				repo: &MockImageRepository{
					GetImageByUserIDandNameFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return &domain.Image{ID: imageID}, nil
					},
					UpdateImageFunc: func(ctx context.Context, id uuid.UUID) error {
						return nil
					},
				},
				storage: &MockStorageService{
					DownloadFunc: func(ctx context.Context, objectName string) ([]byte, error) {
						return mockImageBytes, nil
					},
					UploadFunc: func(ctx context.Context, objectName string, data []byte) error {
						return fmt.Errorf("error uploading image")
					},
				},
				cache: &MockCacheService{
					GetFunc: func(key string) ([]byte, error) {
						return nil, nil
					},
					SetFunc: func(key string, data []byte, duration time.Duration) error {
						return nil
					},
				},
				transformationService: transformations.New(10, 100),
			},
			args: args{
				userID:          validUUID,
				imageName:       validImageName,
				transformations: transformationsList,
			},
			wantErr: true,
		},
		{
			name: "Error setting transformed image in cache",
			fields: fields{
				repo: &MockImageRepository{
					GetImageByUserIDandNameFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return &domain.Image{ID: imageID}, nil
					},
					UpdateImageFunc: func(ctx context.Context, id uuid.UUID) error {
						return nil
					},
				},
				storage: &MockStorageService{
					DownloadFunc: func(ctx context.Context, objectName string) ([]byte, error) {
						return mockImageBytes, nil
					},
					UploadFunc: func(ctx context.Context, objectName string, data []byte) error {
						return nil
					},
				},
				cache: &MockCacheService{
					SetFunc: func(key string, data []byte, duration time.Duration) error {
						return fmt.Errorf("cache error")
					},
					GetFunc: func(key string) ([]byte, error) {
						return nil, nil
					},
				},
				transformationService: transformations.New(10, 100),
			},
			args: args{
				userID:          validUUID,
				imageName:       validImageName,
				transformations: transformationsList,
			},
			wantErr: true,
		},
		{
			name: "Image not found in repository",
			fields: fields{
				repo: &MockImageRepository{
					GetImageByUserIDandNameFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return nil, fmt.Errorf("image not found")
					},
					UpdateImageFunc: func(ctx context.Context, id uuid.UUID) error {
						return nil
					},
				},
				storage: &MockStorageService{
					DownloadFunc: func(ctx context.Context, objectName string) ([]byte, error) {
						return nil, fmt.Errorf("image not found")
					},
					UploadFunc: func(ctx context.Context, objectName string, data []byte) error {
						return nil
					},
				},
				cache: &MockCacheService{
					GetFunc: func(key string) ([]byte, error) {
						return nil, nil
					},
					SetFunc: func(key string, data []byte, duration time.Duration) error {
						return nil
					},
				},
				transformationService: transformations.New(10, 100),
			},
			args: args{
				userID:          validUUID,
				imageName:       validImageName,
				transformations: transformationsList,
			},
			wantErr: true,
		},
		{
			name: "Invalid transformation type",
			fields: fields{
				repo: &MockImageRepository{
					GetImageByUserIDandNameFunc: func(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
						return &domain.Image{ID: imageID}, nil
					},
					UpdateImageFunc: func(ctx context.Context, id uuid.UUID) error {
						return nil
					},
				},
				storage: &MockStorageService{
					DownloadFunc: func(ctx context.Context, objectName string) ([]byte, error) {
						return mockImageBytes, nil
					},
					UploadFunc: func(ctx context.Context, objectName string, data []byte) error {
						return nil
					},
				},
				cache: &MockCacheService{
					GetFunc: func(key string) ([]byte, error) {
						return nil, nil
					},
					SetFunc: func(key string, data []byte, duration time.Duration) error {
						return nil
					},
				},
				transformationService: transformations.New(10, 100),
			},
			args: args{
				userID:    validUUID,
				imageName: validImageName,
				transformations: []domain.Transformation{
					{
						Type:    "UnsupportedTransformation",
						Options: map[string]float64{},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ImageService{
				repo:                  tt.fields.repo,
				storage:               tt.fields.storage,
				cache:                 tt.fields.cache,
				transformationService: tt.fields.transformationService,
			}
			if err := s.ApplyTransformations(tt.args.userID, tt.args.imageName, tt.args.transformations); (err != nil) != tt.wantErr {
				t.Errorf("ApplyTransformations() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func generatePNGImage() []byte {
	// Create a simple 10x10 image with a solid color
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{B: 255, A: 255})
		}
	}

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil
	}

	return buf.Bytes()
}
