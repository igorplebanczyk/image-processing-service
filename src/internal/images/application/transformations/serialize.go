package transformations

import (
	"bytes"
	"fmt"
	"image"
	commonerrors "image-processing-service/src/internal/common/errors"
	"image/jpeg"
	"image/png"
)

func serialize(img image.Image, format string) ([]byte, error) {
	var buf bytes.Buffer

	switch format {
	case "jpeg":
		if err := jpeg.Encode(&buf, img, nil); err != nil {
			return nil, commonerrors.NewInternal("failed to encode image as JPEG")
		}
	case "png":
		if err := png.Encode(&buf, img); err != nil {
			return nil, commonerrors.NewInternal("failed to encode image as PNG")
		}
	default:
		return nil, commonerrors.NewInternal("unsupported image format")
	}

	return buf.Bytes(), nil
}

func deserialize(data []byte) (image.Image, string, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", commonerrors.NewInvalidInput(fmt.Sprintf("error decoding image: %v", err))
	}

	return img, format, nil
}
