package worker

import (
	"github.com/google/uuid"
	"image-processing-service/src/internal/images/domain"
)

func getAllDanglingImagesNames(imagesNamesStorage []string, imagesIDs []uuid.UUID) ([]string, error) {
	var danglingImagesIDs []string
	for _, imageNameStorage := range imagesNamesStorage {
		for _, imageID := range imagesIDs {
			if domain.CreateFullImageObjectName(imageID) == imageNameStorage || domain.CreatePreviewImageObjectName(imageID) == imageNameStorage {
				continue
			}
			danglingImagesIDs = append(danglingImagesIDs, imageNameStorage)
		}
	}

	return danglingImagesIDs, nil
}
