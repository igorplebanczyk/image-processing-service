package transformations

import (
	"image"
	"image/color"
	"reflect"
	"testing"
)

func Test_deserialize(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    image.Image
		want1   string
		wantErr bool
	}{
		{
			name:    "Valid PNG",
			args:    args{data: generateTestSerializedImg("png")},
			want:    generateTestImage(),
			want1:   "png",
			wantErr: false,
		},
		{
			name:    "Invalid image",
			args:    args{data: []byte{}},
			want:    nil,
			want1:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := deserialize(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("deserialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deserialize() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("deserialize() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_serialize(t *testing.T) {
	type args struct {
		img    image.Image
		format string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "Valid PNG",
			args:    args{img: generateTestImage(), format: "png"},
			want:    generateTestSerializedImg("png"),
			wantErr: false,
		},
		{
			name:    "Invalid format",
			args:    args{img: generateTestImage(), format: "invalid"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := serialize(tt.args.img, tt.args.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("serialize() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func generateTestImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	img.Set(1, 0, color.RGBA{G: 255, A: 255})
	img.Set(0, 1, color.RGBA{B: 255, A: 255})
	img.Set(1, 1, color.RGBA{R: 255, G: 255, A: 255})
	return img
}

func generateTestSerializedImg(format string) []byte {
	bytes, _ := serialize(generateTestImage(), format)
	return bytes
}
