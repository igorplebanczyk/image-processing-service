package worker

import "fmt"

func getAllDanglingImagesNames(imagesNamesStorage []string, imagesDB []image) ([]string, error) {
	var danglingImages []string
	for _, imageName := range imagesNamesStorage {
		for _, imageDB := range imagesDB {
			if imageName == fmt.Sprintf("%s-%s", imageDB.userID, imageDB.name) {
				continue
			}
			danglingImages = append(danglingImages, imageName)
		}
	}

	return danglingImages, nil
}
