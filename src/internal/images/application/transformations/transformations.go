package transformations

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"math"
)

func resize(img image.Image, options map[string]float64) (image.Image, error) {
	width, ok := options["width"]
	if !ok {
		return nil, fmt.Errorf("resize option 'width' is required and must be a number")
	}

	height, ok := options["height"]
	if !ok {
		return nil, fmt.Errorf("resize option 'height' is required and must be a number")
	}

	return imaging.Resize(img, int(width), int(height), imaging.Lanczos), nil
}

func crop(img image.Image, options map[string]float64) (image.Image, error) {
	width, ok := options["width"]
	if !ok {
		return nil, fmt.Errorf("crop option 'width' is required and must be a number")
	}

	height, ok := options["height"]
	if !ok {
		return nil, fmt.Errorf("crop option 'height' is required and must be a number")
	}

	return imaging.CropCenter(img, int(width), int(height)), nil
}

func rotate(img image.Image, options map[string]float64) (image.Image, error) {
	angle, ok := options["angle"]
	if !ok {
		return nil, fmt.Errorf("rotate option 'angle' is required and must be a number")
	}

	return imaging.Rotate(img, angle, color.Transparent), nil
}

func grayscale(img image.Image) image.Image {
	return imaging.Grayscale(img)
}

func invert(img image.Image) image.Image {
	return imaging.Invert(img)
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

func adjustBrightness(img image.Image, options map[string]float64) (image.Image, error) {
	factor, ok := options["factor"]
	if !ok {
		return nil, fmt.Errorf("adjust_brightness option 'factor' is required and must be a number")
	}

	return imaging.AdjustBrightness(img, factor), nil
}

func adjustContrast(img image.Image, options map[string]float64) (image.Image, error) {
	factor, ok := options["factor"]
	if !ok {
		return nil, fmt.Errorf("adjust_contrast option 'factor' is required and must be a number")
	}

	return imaging.AdjustContrast(img, factor), nil
}

func adjustSaturation(img image.Image, options map[string]float64) (image.Image, error) {
	factor, ok := options["factor"]
	if !ok {
		return nil, fmt.Errorf("adjust_saturation option 'factor' is required and must be a number")
	}

	return imaging.AdjustSaturation(img, factor), nil
}

func blur(img image.Image, options map[string]float64) (image.Image, error) {
	factor, ok := options["factor"]
	if !ok {
		return nil, fmt.Errorf("blur option 'sigma' is required and must be a number")
	}

	return imaging.Blur(img, factor), nil
}

func sharpen(img image.Image, options map[string]float64) (image.Image, error) {
	factor, ok := options["factor"]
	if !ok {
		return nil, fmt.Errorf("sharpen option 'sigma' is required and must be a number")
	}

	return imaging.Sharpen(img, factor), nil
}
