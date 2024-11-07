package transformation

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"math"
)

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

func (s *Service) applyFilter(img image.Image, options map[string]any) (image.Image, error) {
	filterType, ok := options["filter"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'filter' in options")
	}

	switch filterType {
	case "grayscale":
		return grayscale(img), nil
	case "sepia":
		return sepia(img), nil
	default:
		return nil, fmt.Errorf("unknown filter type: %s", filterType)
	}
}

func grayscale(img image.Image) image.Image {
	return imaging.Grayscale(img)
}

func sepia(img image.Image) image.Image {
	return imaging.AdjustFunc(img, func(c color.NRGBA) color.NRGBA {
		tr := float64(c.R)*0.393 + float64(c.G)*0.769 + float64(c.B)*0.189
		tg := float64(c.R)*0.349 + float64(c.G)*0.686 + float64(c.B)*0.168
		tb := float64(c.R)*0.272 + float64(c.G)*0.534 + float64(c.B)*0.131
		return color.NRGBA{
			R: uint8(math.Min(tr, 255)),
			G: uint8(math.Min(tg, 255)),
			B: uint8(math.Min(tb, 255)),
			A: c.A,
		}
	})
}
