package domain

type Transformation struct {
	Type    TransformationType
	Options map[string]float64
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
