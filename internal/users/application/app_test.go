package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/internal/users/domain"
	"reflect"
	"testing"
	"time"
)

// Mocks

type MockUserRepository struct {
	CreateUserFunc  func(ctx context.Context, username, email, password string) (*domain.User, error)
	GetUserByIDFunc func(ctx context.Context, id uuid.UUID) (*domain.User, error)
	UpdateUserFunc  func(ctx context.Context, id uuid.UUID, username, email string) error
	DeleteUserFunc  func(ctx context.Context, id uuid.UUID) error
}

func (m *MockUserRepository) CreateUser(ctx context.Context, username, email, password string) (*domain.User, error) {
	return m.CreateUserFunc(ctx, username, email, password)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return m.GetUserByIDFunc(ctx, id)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, id uuid.UUID, username, email string) error {
	return m.UpdateUserFunc(ctx, id, username, email)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return m.DeleteUserFunc(ctx, id)
}

// Tests

func TestUserService_DeleteUser(t *testing.T) {
	type fields struct {
		repo domain.UserRepository
	}
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful delete",
			fields: fields{
				repo: &MockUserRepository{
					DeleteUserFunc: func(ctx context.Context, id uuid.UUID) error {
						return nil // simulate successful deletion
					},
				},
			},
			args:    args{userID: uuid.New()},
			wantErr: false,
		},
		{
			name: "failed delete",
			fields: fields{
				repo: &MockUserRepository{
					DeleteUserFunc: func(ctx context.Context, id uuid.UUID) error {
						return fmt.Errorf("delete failed") // simulate failure
					},
				},
			},
			args:    args{userID: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserService{repo: tt.fields.repo}
			if err := s.DeleteUser(tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	testID := uuid.New()
	type fields struct {
		repo domain.UserRepository
	}
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.User
		wantErr bool
	}{
		{
			name: "successful get user",
			fields: fields{
				repo: &MockUserRepository{
					GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
						return &domain.User{ID: id, Username: "testuser", Email: "test@example.com"}, nil
					},
				},
			},
			args:    args{userID: testID},
			want:    &domain.User{ID: testID, Username: "testuser", Email: "test@example.com"},
			wantErr: false,
		},
		{
			name: "failed to get user",
			fields: fields{
				repo: &MockUserRepository{
					GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
						return nil, fmt.Errorf("user not found")
					},
				},
			},
			args:    args{userID: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserService{repo: tt.fields.repo}
			got, err := s.GetUser(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_Register(t *testing.T) {
	testID := uuid.New()
	testTime := time.Now()

	type fields struct {
		repo domain.UserRepository
	}
	type args struct {
		username string
		email    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.User
		wantErr bool
	}{
		{
			name: "successful registration",
			fields: fields{
				repo: &MockUserRepository{
					CreateUserFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
						return &domain.User{ID: testID, Username: username, Email: email, CreatedAt: testTime, UpdatedAt: testTime}, nil
					},
				},
			},
			args:    args{username: "testuser", email: "test@example.com", password: "password123A!"},
			want:    &domain.User{ID: testID, Username: "testuser", Email: "test@example.com", CreatedAt: testTime, UpdatedAt: testTime},
			wantErr: false,
		},
		{
			name: "failed registration",
			fields: fields{
				repo: &MockUserRepository{
					CreateUserFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
						return nil, fmt.Errorf("registration failed")
					},
				},
			},
			args:    args{username: "testuser", email: "test@example.com", password: "password123"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserService{repo: tt.fields.repo}
			got, err := s.Register(tt.args.username, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Register() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	type fields struct {
		repo domain.UserRepository
	}
	type args struct {
		userID   uuid.UUID
		username string
		email    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful update",
			fields: fields{
				repo: &MockUserRepository{
					GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
						return &domain.User{ID: id, Username: "testuser", Email: "test@example.com"}, nil
					},
					UpdateUserFunc: func(ctx context.Context, id uuid.UUID, username, email string) error {
						return nil
					},
				},
			},
			args:    args{userID: uuid.New(), username: "newuser", email: "newemail@example.com"},
			wantErr: false,
		},
		{
			name: "failed update - user not found",
			fields: fields{
				repo: &MockUserRepository{
					GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
						return nil, fmt.Errorf("user not found")
					},
				},
			},
			args:    args{userID: uuid.New(), username: "newuser", email: "newemail@example.com"},
			wantErr: true,
		},
		{
			name: "failed update - validation error",
			fields: fields{
				repo: &MockUserRepository{
					GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
						return &domain.User{ID: id, Username: "testuser", Email: "test@example.com"}, nil
					},
				},
			},
			args:    args{userID: uuid.New(), username: "", email: "invalid-email"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserService{repo: tt.fields.repo}
			if err := s.UpdateUser(tt.args.userID, tt.args.username, tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
