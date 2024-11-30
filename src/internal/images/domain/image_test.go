package domain

import (
	"encoding/base64"
	"github.com/google/uuid"
	"testing"
)

func TestValidateName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Valid name",
			args{name: "JohnDoe"},
			false,
		},
		{
			"Name too short",
			args{name: "JD"},
			true,
		},
		{
			"Name too long",
			args{name: string(make([]byte, 129))},
			true,
		},
		{
			"Name contains spaces",
			args{name: "John Doe"},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
	type args struct {
		description string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Valid description",
			args{description: "This is a description"},
			false,
		},
		{
			"Description too long",
			args{description: string(make([]byte, 1025))},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDescription(tt.args.description)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateDescription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateImage(t *testing.T) {
	type args struct {
		imageBytes []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Valid Image",
			args{imageBytes: validImage},
			false,
		},
		{
			"Image too large",
			args{imageBytes: make([]byte, MaxImageSize+1)},
			true,
		},
		{
			"Invalid Image Format",
			args{imageBytes: []byte("invalid")},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateImage(tt.args.imageBytes)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateRawImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDetermineImageMetadataToUpdate(t *testing.T) {
	type args struct {
		existingImageMetadata *ImageMetadata
		newName               string
		newDescription        string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			"Update name and description",
			args{
				existingImageMetadata: &ImageMetadata{
					Name:        "OldName",
					Description: "OldDescription",
				},
				newName:        "NewName",
				newDescription: "NewDescription",
			},
			"NewName",
			"NewDescription",
			false,
		},
		{
			"Update name",
			args{
				existingImageMetadata: &ImageMetadata{
					Name:        "OldName",
					Description: "OldDescription",
				},
				newName:        "NewName",
				newDescription: "",
			},
			"NewName",
			"OldDescription",
			false,
		},
		{
			"Update description",
			args{
				existingImageMetadata: &ImageMetadata{
					Name:        "OldName",
					Description: "OldDescription",
				},
				newName:        "",
				newDescription: "NewDescription",
			},
			"OldName",
			"NewDescription",
			false,
		},
		{
			"Update nothing",
			args{
				existingImageMetadata: &ImageMetadata{
					Name:        "OldName",
					Description: "OldDescription",
				},
				newName:        "",
				newDescription: "",
			},
			"",
			"",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := DetermineImageMetadataToUpdate(tt.args.existingImageMetadata, tt.args.newName, tt.args.newDescription)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetermineImageMetadataToUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DetermineImageMetadataToUpdate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DetermineImageMetadataToUpdate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCreateFullImageObjectName(t *testing.T) {
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Create full image object name",
			args{id: uuid.MustParse("00000000-0000-0000-0000-000000000000")},
			"full-00000000-0000-0000-0000-000000000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateFullImageObjectName(tt.args.id); got != tt.want {
				t.Errorf("CreateFullImageObjectName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreatePreviewImageObjectName(t *testing.T) {
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Create preview image object name",
			args{id: uuid.MustParse("00000000-0000-0000-0000-000000000000")},
			"prev-00000000-0000-0000-0000-000000000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreatePreviewImageObjectName(tt.args.id); got != tt.want {
				t.Errorf("CreatePreviewImageObjectName() = %v, want %v", got, tt.want)
			}
		})
	}
}

var validImage, _ = base64.StdEncoding.DecodeString(`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`)
