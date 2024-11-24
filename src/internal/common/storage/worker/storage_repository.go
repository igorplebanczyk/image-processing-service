package worker

import (
	"context"
	"fmt"
	"image-processing-service/src/internal/common/storage"
	"time"
)

type imagesStorageRepository struct {
	storage *storage.Service
}

func newImagesStorageRepository(storage *storage.Service) imagesStorageRepository {
	return imagesStorageRepository{
		storage: storage,
	}
}

func (r *imagesStorageRepository) getAllImagesNames() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pager := r.storage.Client().NewListBlobsFlatPager(r.storage.ContainerName(), nil)

	var imagesNames []string
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, blob := range resp.Segment.BlobItems {
			imagesNames = append(imagesNames, *blob.Name)
		}
	}

	return imagesNames, nil
}

func (r *imagesStorageRepository) deleteImage(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := r.storage.Client().DeleteBlob(ctx, r.storage.ContainerName(), name, nil)
	if err != nil {
		return fmt.Errorf("failed to delete image %s: %w", name, err)
	}

	return nil
}
