package application

import "testing"

func Test_hashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Valid password",
			args:    args{password: "Password123!"},
			wantErr: false,
		},
		{
			name:    "Password too long",
			args:    args{password: "EfkUEP2pZJjBdgOzPzFqpXghkT0N8pcwmonZyNnXtkLfBD7tviIoVCLBt0HSVNN06RKqlUIAA"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := hashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("hashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
