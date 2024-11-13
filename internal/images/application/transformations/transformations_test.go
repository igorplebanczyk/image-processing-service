package transformations

import (
	"image"
	"image/color"
	"testing"
)

func TestSerialize(t *testing.T) {
	// Create a simple image (a 2x2 red image)
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	img.Set(1, 0, color.RGBA{G: 255, A: 255})
	img.Set(0, 1, color.RGBA{B: 255, A: 255})
	img.Set(1, 1, color.RGBA{R: 255, G: 255, A: 255})

	tests := []struct {
		format  string
		wantErr bool
	}{
		{"jpeg", false},
		{"png", false},
		{"bmp", true}, // unsupported format
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			_, err := serialize(img, tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("serialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeserialize(t *testing.T) {
	// Create a simple image (a 2x2 red image)
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	img.Set(1, 0, color.RGBA{G: 255, A: 255})
	img.Set(0, 1, color.RGBA{B: 255, A: 255})
	img.Set(1, 1, color.RGBA{R: 255, G: 255, A: 255})

	// Serialize the image into PNG format
	data, err := serialize(img, "png")
	if err != nil {
		t.Fatalf("serialize() failed: %v", err)
	}

	tests := []struct {
		data    []byte
		wantErr bool
	}{
		{data, false},                  // valid PNG data
		{[]byte("invalid data"), true}, // invalid data
	}

	for _, tt := range tests {
		t.Run("Deserialize", func(t *testing.T) {
			decodedImg, format, err := deserialize(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("deserialize() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && decodedImg == nil {
				t.Errorf("expected a decoded image, got nil")
			}
			if !tt.wantErr && format != "png" {
				t.Errorf("expected format 'png', got '%s'", format)
			}
		})
	}
}
