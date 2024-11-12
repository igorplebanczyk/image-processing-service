package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/internal/users/domain"
	"time"
)

type UserService struct {
	repo domain.UserRepository
}

func NewService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(username, email, password string) (*domain.User, error) {
	err := domain.ValidateUsername(username)
	if err != nil {
		return nil, fmt.Errorf("invalid username: %w", err)
	}

	err = domain.ValidateEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	err = domain.ValidatePassword(password)
	if err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	hashedPassword, err := domain.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.repo.CreateUser(ctx, username, email, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUser(userID uuid.UUID) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

func (s *UserService) UpdateUser(userID uuid.UUID, username, email string) error {
	err := domain.ValidateUsername(username)
	if err != nil {
		return fmt.Errorf("invalid username: %w", err)
	}

	err = domain.ValidateEmail(email)
	if err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = s.repo.UpdateUser(ctx, userID, username, email)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

func (s *UserService) DeleteUser(userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.repo.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}
