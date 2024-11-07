package transformation

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
)

func serialize(img image.Image, format string) ([]byte, error) {
	var buf bytes.Buffer

	switch format {
	case "jpeg":
		if err := jpeg.Encode(&buf, img, nil); err != nil {
			return nil, fmt.Errorf("failed to encode image as JPEG: %w", err)
		}
	case "png":
		if err := png.Encode(&buf, img); err != nil {
			return nil, fmt.Errorf("failed to encode image as PNG: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	return buf.Bytes(), nil
}

func deserialize(data []byte) (image.Image, string, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", fmt.Errorf("error decoding image: %w", err)
	}

	return img, format, nil
}
