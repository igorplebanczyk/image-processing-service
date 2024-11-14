package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/internal/users/domain"
	"testing"
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

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		email         string
		password      string
		expectedError error
		mockFunc      func(ctx context.Context, username, email, password string) (*domain.User, error)
	}{
		{
			name:     "successful registration",
			username: "new_user",
			email:    "new_user@example.com",
			password: "ValidPassword123!",
			mockFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
				return &domain.User{
					ID:       uuid.New(),
					Username: "new_user",
					Email:    "new_user@example.com",
					Password: "ValidPassword123",
				}, nil
			},
		},
		{
			name:          "invalid username",
			username:      "",
			email:         "valid@example.com",
			password:      "ValidPassword123!",
			expectedError: fmt.Errorf("invalid username: username cannot be empty"),
			mockFunc:      nil,
		},
		{
			name:          "invalid email",
			username:      "new_user",
			email:         "invalid-email",
			password:      "ValidPassword123!",
			expectedError: fmt.Errorf("invalid email: mail: missing '@' or angle-addr"),
			mockFunc:      nil,
		},
		{
			name:          "invalid password",
			username:      "new_user",
			email:         "new_user@example.com",
			password:      "short",
			expectedError: fmt.Errorf("invalid password: password must be between 8 and 32 characters"),
			mockFunc:      nil,
		},
		{
			name:          "error in user creation",
			username:      "new_user",
			email:         "new_user@example.com",
			password:      "ValidPassword123!",
			expectedError: fmt.Errorf("error creating user: %w", errors.New("mock error")),
			mockFunc: func(ctx context.Context, username, email, password string) (*domain.User, error) {
				return nil, errors.New("mock error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &MockUserRepository{
				CreateUserFunc: tt.mockFunc,
			}

			userService := &UserService{
				repo: mockUserRepo,
			}

			user, err := userService.Register(tt.username, tt.email, tt.password)

			if tt.expectedError != nil && err == nil {
				t.Errorf("expected error %v, got nil", tt.expectedError)
			} else if tt.expectedError == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			} else if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if tt.expectedError == nil && user == nil {
				t.Error("expected user to be returned, got nil")
			}

			if user != nil && user.Username != tt.username {
				t.Errorf("expected username %v, got %v", tt.username, user.Username)
			}
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        uuid.UUID
		mockFunc      func(ctx context.Context, id uuid.UUID) (*domain.User, error)
		expectedUser  *domain.User
		expectedError error
	}{
		{
			name:   "successful user retrieval",
			userID: uuid.New(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				return &domain.User{
					ID:       id,
					Username: "existing_user",
					Email:    "existing_user@example.com",
					Password: "ValidPassword123",
				}, nil
			},
			expectedUser: &domain.User{
				Username: "existing_user",
				Email:    "existing_user@example.com",
				Password: "ValidPassword123",
			},
			expectedError: nil,
		},
		{
			name:   "user not found",
			userID: uuid.New(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				return nil, errors.New("user not found")
			},
			expectedUser:  nil,
			expectedError: fmt.Errorf("error getting user: %w", errors.New("user not found")),
		},
		{
			name:   "database error",
			userID: uuid.New(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				return nil, errors.New("database connection error")
			},
			expectedUser:  nil,
			expectedError: fmt.Errorf("error getting user: %w", errors.New("database connection error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &MockUserRepository{
				GetUserByIDFunc: tt.mockFunc,
			}

			userService := &UserService{
				repo: mockUserRepo,
			}

			user, err := userService.GetUser(tt.userID)

			if tt.expectedError != nil && err == nil {
				t.Errorf("expected error %v, got nil", tt.expectedError)
			} else if tt.expectedError == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			} else if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if tt.expectedError == nil && user == nil {
				t.Error("expected user to be returned, got nil")
			}

			if tt.expectedUser != nil && user != nil {
				if user.Username != tt.expectedUser.Username {
					t.Errorf("expected username %v, got %v", tt.expectedUser.Username, user.Username)
				}
				if user.Email != tt.expectedUser.Email {
					t.Errorf("expected email %v, got %v", tt.expectedUser.Email, user.Email)
				}
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         uuid.UUID
		username       string
		email          string
		mockGetUser    func(ctx context.Context, id uuid.UUID) (*domain.User, error)
		mockUpdateUser func(ctx context.Context, id uuid.UUID, username, email string) error
		expectedError  error
	}{
		{
			name:     "successful update",
			userID:   uuid.New(),
			username: "updated_user",
			email:    "updated_user@example.com",
			mockGetUser: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				return &domain.User{
					ID:       id,
					Username: "existing_user",
					Email:    "existing_user@example.com",
				}, nil
			},
			mockUpdateUser: func(ctx context.Context, id uuid.UUID, username, email string) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name:     "no data to update",
			userID:   uuid.New(),
			username: "",
			email:    "",
			mockGetUser: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				return &domain.User{
					ID:       id,
					Username: "existing_user",
					Email:    "existing_user@example.com",
				}, nil
			},
			mockUpdateUser: nil,
			expectedError:  fmt.Errorf("no data to update"),
		},
		{
			name:     "error getting user",
			userID:   uuid.New(),
			username: "updated_user",
			email:    "updated_user@example.com",
			mockGetUser: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				return nil, errors.New("user not found")
			},
			mockUpdateUser: nil,
			expectedError:  fmt.Errorf("error getting user: %w", errors.New("user not found")),
		},
		{
			name:     "error updating user",
			userID:   uuid.New(),
			username: "updated_user",
			email:    "updated_user@example.com",
			mockGetUser: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
				return &domain.User{
					ID:       id,
					Username: "existing_user",
					Email:    "existing_user@example.com",
				}, nil
			},
			mockUpdateUser: func(ctx context.Context, id uuid.UUID, username, email string) error {
				return errors.New("update failed")
			},
			expectedError: fmt.Errorf("error updating user: %w", errors.New("update failed")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &MockUserRepository{
				GetUserByIDFunc: tt.mockGetUser,
				UpdateUserFunc:  tt.mockUpdateUser,
			}

			userService := &UserService{
				repo: mockUserRepo,
			}

			err := userService.UpdateUser(tt.userID, tt.username, tt.email)

			if tt.expectedError != nil && err == nil {
				t.Errorf("expected error %v, got nil", tt.expectedError)
			} else if tt.expectedError == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			} else if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        uuid.UUID
		mockFunc      func(ctx context.Context, id uuid.UUID) error
		expectedError error
	}{
		{
			name:   "successful deletion",
			userID: uuid.New(),
			mockFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
		},
		{
			name:   "error deleting user",
			userID: uuid.New(),
			mockFunc: func(ctx context.Context, id uuid.UUID) error {
				return errors.New("mock error")
			},
			expectedError: fmt.Errorf("error deleting user: %w", errors.New("mock error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &MockUserRepository{
				DeleteUserFunc: tt.mockFunc,
			}

			userService := &UserService{
				repo: mockUserRepo,
			}

			err := userService.DeleteUser(tt.userID)

			// Check for errors
			if tt.expectedError != nil && err == nil {
				t.Errorf("expected error %v, got nil", tt.expectedError)
			} else if tt.expectedError == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			} else if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}
