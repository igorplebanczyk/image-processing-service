package transformation

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"math"
)

func resize(img image.Image, options map[string]any) (image.Image, error) {
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

func crop(img image.Image, options map[string]any) (image.Image, error) {
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

func rotate(img image.Image, options map[string]any) (image.Image, error) {
	angle, ok := options["angle"].(float64)
	if !ok {
		return nil, fmt.Errorf("rotate option 'angle' is required and must be a number")
	}

	return imaging.Rotate(img, angle, color.Transparent), nil
}

func applyFilter(img image.Image, options map[string]any) (image.Image, error) {
	filterType, ok := options["filter"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'filter' in options")
	}

	switch filterType {
	case "grayscale":
		return imaging.Grayscale(img), nil
	case "sepia":
		return sepia(img), nil
	case "invert":
		return imaging.Invert(img), nil
	default:
		return nil, fmt.Errorf("unknown filter type: %s", filterType)
	}
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

func adjust(img image.Image, options map[string]any) (image.Image, error) {
	adjustType, ok := options["adjust"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'adjust' in options")
	}

	factor, ok := options["factor"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'factor' in options")
	}

	switch adjustType {
	case "brightness":
		return imaging.AdjustBrightness(img, factor), nil
	case "contrast":
		return imaging.AdjustContrast(img, factor), nil
	case "saturation":
		return imaging.AdjustSaturation(img, factor), nil
	case "blur":
		return imaging.Blur(img, factor), nil
	case "sharpen":
		return imaging.Sharpen(img, factor), nil
	default:
		return nil, fmt.Errorf("unknown adjust type: %s", adjustType)
	}
}
