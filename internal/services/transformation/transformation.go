package transformation

import (
	"fmt"
	"image"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

const (
	resize      = "resize"
	crop        = "crop"
	rotate      = "rotate"
	applyFilter = "apply_filter"
)

func (s *Service) Transform(imageBytes []byte, transformation string, options map[string]any) ([]byte, error) {
	img, format, err := deserialize(imageBytes)
	if err != nil {
		return nil, fmt.Errorf("error deserializing image: %w", err)
	}

	var transformed image.Image

	switch transformation {
	case resize:
		transformed, err = s.resize(img, options)
	case crop:
		transformed, err = s.crop(img, options)
	case rotate:
		transformed, err = s.rotate(img, options)
	case applyFilter:
		transformed, err = s.applyFilter(img, options)
	default:
		return nil, fmt.Errorf("unknown transformation: %s", transformation)
	}
	if err != nil {
		return nil, fmt.Errorf("error during %s transformation: %w", transformation, err)
	}

	return serialize(transformed, format)
}
