package domain

import (
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

func TestValidateRawImage(t *testing.T) {
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
			args{imageBytes: make([]byte, MaxImageSize-1)},
			false,
		},
		{
			"Image too large",
			args{imageBytes: make([]byte, MaxImageSize+1)},
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

func TestCreateObjectName(t *testing.T) {
	testUUID := uuid.New()

	type args struct {
		userID    uuid.UUID
		imageName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Valid Object Name",
			args{userID: testUUID, imageName: "image"},
			testUUID.String() + "-image",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateFullImageObjectName(tt.args.userID, tt.args.imageName)
			if got != tt.want {
				t.Fatalf("CreateObjectName() = %v, want %v", got, tt.want)
			}
		})
	}
}
