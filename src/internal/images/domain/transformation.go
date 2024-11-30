package domain

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
