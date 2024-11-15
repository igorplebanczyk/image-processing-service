package transformations

import (
	"errors"
	"fmt"
	"image"
	"image-processing-service/internal/images/domain"
)

type task struct {
	imageBytes     []byte
	transformation domain.TransformationType
	options        map[string]float64
	result         chan []byte
	err            chan error
}

func (s *Service) startWorkerPool() {
	for i := 0; i < s.workerCount; i++ {
		s.wg.Add(1)
		go s.worker()
	}
}

func (s *Service) worker() {
	defer s.wg.Done()
	for task := range s.taskQueue {
		result, err := s.processTask(task.imageBytes, task.transformation, task.options)
		task.result <- result
		task.err <- err
		close(task.result)
		close(task.err)
	}
}

func (s *Service) processTask(imageBytes []byte, transformation domain.TransformationType, options map[string]float64) ([]byte, error) {
	img, format, err := deserialize(imageBytes)
	if err != nil {
		return nil, fmt.Errorf("error deserializing image: %w", err)
	}

	var transformed image.Image
	switch transformation {
	case domain.Resize:
		transformed, err = resize(img, options)
	case domain.Crop:
		transformed, err = crop(img, options)
	case domain.Rotate:
		transformed, err = rotate(img, options)
	case domain.Grayscale:
		transformed = grayscale(img)
	case domain.Invert:
		transformed = invert(img)
	case domain.Sepia:
		transformed = sepia(img)
	case domain.AdjustBrightness:
		transformed, err = adjustBrightness(img, options)
	case domain.AdjustContrast:
		transformed, err = adjustContrast(img, options)
	case domain.AdjustSaturation:
		transformed, err = adjustSaturation(img, options)
	case domain.Blur:
		transformed, err = blur(img, options)
	case domain.Sharpen:
		transformed, err = sharpen(img, options)
	default:
		return nil, errors.Join(domain.ErrInvalidRequest, fmt.Errorf("unsupported transformation: %s", transformation))
	}
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("error applying transformation: %w", err))
	}

	return serialize(transformed, format)
}
