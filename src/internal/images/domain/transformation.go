package domain

import "fmt"

type Transformation struct {
	Type    TransformationType
	Options map[TransformationOptionType]float64
}

type TransformationType string

const (
	Resize           TransformationType = "resize"
	Crop             TransformationType = "crop"
	Rotate           TransformationType = "rotate"
	Grayscale        TransformationType = "grayscale"
	Sepia            TransformationType = "sepia"
	Invert           TransformationType = "invert"
	AdjustBrightness TransformationType = "adjust_brightness"
	AdjustContrast   TransformationType = "adjust_contrast"
	AdjustSaturation TransformationType = "adjust_saturation"
	Blur             TransformationType = "blur"
	Sharpen          TransformationType = "sharpen"
)

type TransformationOptionType string

const (
	Width  TransformationOptionType = "width"
	Height TransformationOptionType = "height"
	Angle  TransformationOptionType = "angle"
	Factor TransformationOptionType = "factor"
)

func NewTransformation(transformationType TransformationType, options map[TransformationOptionType]float64) *Transformation {
	return &Transformation{
		Type:    transformationType,
		Options: options,
	}
}

func ValidateTransformation(transformation *Transformation) error {
	switch transformation.Type {
	case Resize, Crop:
		_, widthOk := transformation.Options[Width]
		_, heightOk := transformation.Options[Height]
		if !widthOk || !heightOk ||
			transformation.Options[Width] <= 0 || transformation.Options[Height] <= 0 ||
			transformation.Options[Width] > 4096 || transformation.Options[Height] > 4096 {
			return fmt.Errorf("invalid width or height")
		}
	case Rotate:
		_, angleOk := transformation.Options[Angle]
		if !angleOk || transformation.Options[Angle] < 0 || transformation.Options[Angle] > 360 {
			return fmt.Errorf("invalid angle")
		}
	case Grayscale, Sepia, Invert:
		if len(transformation.Options) > 0 {
			return fmt.Errorf("unexpected options")
		}
	case AdjustBrightness, AdjustContrast, AdjustSaturation, Blur, Sharpen:
		_, factorOk := transformation.Options[Factor]
		if !factorOk {
			return fmt.Errorf("missing factor")
		}
	default:
		return fmt.Errorf("unknown transformation type")
	}

	return nil
}
