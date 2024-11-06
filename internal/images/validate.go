package images

import (
	"fmt"
	"github.com/google/uuid"
)

func validate(r Repository, userID uuid.UUID, imageName string) (bool, error) {
	if imageName == "" {
		return false, nil
	}

	userImages, err := r.GetImagesByUserID(userID)
	if err != nil {
		return false, fmt.Errorf("error getting images by user id: %w", err)
	}

	for _, image := range userImages {
		if image.Name == imageName {
			return false, nil
		}
	}

	return true, nil
}
