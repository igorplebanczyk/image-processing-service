package domain

import (
	"testing"
)

func TestValidateUsername(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Valid username",
			args:    args{username: "username"},
			wantErr: false,
		},
		{
			name:    "Empty username",
			args:    args{username: ""},
			wantErr: true,
		},
		{
			name:    "Username too short",
			args:    args{username: "us"},
			wantErr: true,
		},
		{
			name:    "Username too long",
			args:    args{username: "thisusernameiswaytoolongandshouldefinitelydnotbeaccepted"},
			wantErr: true,
		},
		{
			name:    "Username contains spaces",
			args:    args{username: "user name"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateUsername(tt.args.username); (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Valid email",
			args:    args{email: "valid@example.com"},
			wantErr: false,
		},
		{
			name:    "Empty email",
			args:    args{email: ""},
			wantErr: true,
		},
		{
			name:    "No @ symbol",
			args:    args{email: "invalid"},
			wantErr: true,
		},
		{
			name:    "No domain",
			args:    args{email: "invalid@"},
			wantErr: true,
		},
		{
			name:    "No username",
			args:    args{email: "@invalid.com"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateEmail(tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
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
			name:    "Empty password",
			args:    args{password: ""},
			wantErr: true,
		},
		{
			name:    "Password too short",
			args:    args{password: "Pas123!"},
			wantErr: true,
		},
		{
			name:    "Password too long",
			args:    args{password: "ThisPasswordIsWayTooLongAndShouldDefinitelyNotBeAccepted123!"},
			wantErr: true,
		},
		{
			name:    "No uppercase letter",
			args:    args{password: "password123!"},
			wantErr: true,
		},
		{
			name:    "No lowercase letter",
			args:    args{password: "PASSWORD123!"},
			wantErr: true,
		},
		{
			name:    "No number",
			args:    args{password: "Password!"},
			wantErr: true,
		},
		{
			name:    "No special character",
			args:    args{password: "Password123"},
			wantErr: true,
		},
		{
			name:    "Password contains spaces",
			args:    args{password: "Password 123!"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidatePassword(tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDetermineUserDetailsToUpdate(t *testing.T) {
	type args struct {
		existingUser *User
		newUsername  string
		newEmail     string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "Both username and email need to be updated",
			args: args{
				existingUser: NewUser("oldUsername", "old@example.com", "password123", ""),
				newUsername:  "newUsername",
				newEmail:     "new@example.com",
			},
			want:    "newUsername",
			want1:   "new@example.com",
			wantErr: false,
		},
		{
			name: "Only username needs to be updated",
			args: args{
				existingUser: NewUser("oldUsername", "old@example.com", "password123", ""),
				newUsername:  "newUsername",
				newEmail:     "",
			},
			want:    "newUsername",
			want1:   "old@example.com",
			wantErr: false,
		},
		{
			name: "Only email needs to be updated",
			args: args{
				existingUser: NewUser("oldUsername", "old@example.com", "password123", ""),
				newUsername:  "",
				newEmail:     "new@example.com",
			},
			want:    "oldUsername",
			want1:   "new@example.com",
			wantErr: false,
		},
		{
			name: "No values provided",
			args: args{
				existingUser: NewUser("oldUsername", "old@example.com", "password123", ""),
				newUsername:  "",
				newEmail:     "",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := DetermineUserDetailsToUpdate(tt.args.existingUser, tt.args.newUsername, tt.args.newEmail)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetermineUserDetailsToUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DetermineUserDetailsToUpdate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DetermineUserDetailsToUpdate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
