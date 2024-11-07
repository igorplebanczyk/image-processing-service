package transformation

import (
	"fmt"
	"image"
	"image/color"

	"github.com/disintegration/imaging"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

const (
	resize = "resize"
	crop   = "crop"
	rotate = "rotate"
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
	default:
		return nil, fmt.Errorf("unknown transformation: %s", transformation)
	}
	if err != nil {
		return nil, fmt.Errorf("error during %s transformation: %w", transformation, err)
	}

	return serialize(transformed, format)
}

func (s *Service) resize(img image.Image, options map[string]any) (image.Image, error) {
	width, ok := options["width"].(float64)
	if !ok {
		return nil, fmt.Errorf("resize option 'width' is required and must be a number")
	}

	height, ok := options["height"].(float64)
	if !ok {
		return nil, fmt.Errorf("resize option 'height' is required and must be a number")
	}

	return imaging.Resize(img, int(width), int(height), imaging.Lanczos), nil
}

func (s *Service) crop(img image.Image, options map[string]any) (image.Image, error) {
	width, ok := options["width"].(float64)
	if !ok {
		return nil, fmt.Errorf("crop option 'width' is required and must be a number")
	}

	height, ok := options["height"].(float64)
	if !ok {
		return nil, fmt.Errorf("crop option 'height' is required and must be a number")
	}

	return imaging.CropCenter(img, int(width), int(height)), nil
}

func (s *Service) rotate(img image.Image, options map[string]any) (image.Image, error) {
	angle, ok := options["angle"].(float64)
	if !ok {
		return nil, fmt.Errorf("rotate option 'angle' is required and must be a number")
	}

	return imaging.Rotate(img, angle, color.Transparent), nil
}
