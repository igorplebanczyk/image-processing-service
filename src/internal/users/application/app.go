package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	commonerrors "image-processing-service/src/internal/common/errors"
	"image-processing-service/src/internal/users/domain"
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
		return nil, commonerrors.NewInvalidInput(fmt.Sprintf("invalid username: %v", err))
	}

	err = domain.ValidateEmail(email)
	if err != nil {
		return nil, commonerrors.NewInvalidInput(fmt.Sprintf("invalid email: %v", err))
	}

	err = domain.ValidatePassword(password)
	if err != nil {
		return nil, commonerrors.NewInvalidInput(fmt.Sprintf("invalid password: %v", err))
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, commonerrors.NewInternal(fmt.Sprintf("failed to hash password: %v", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.repo.CreateUser(ctx, username, email, hashedPassword)
	if err != nil {
		return nil, commonerrors.NewInternal(fmt.Sprintf("failed to create user: %v", err))
	}

	return user, nil
}

func (s *UserService) GetUser(userID uuid.UUID) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, commonerrors.NewInternal(fmt.Sprintf("failed to get user from database: %v", err))
	}

	return user, nil
}

func (s *UserService) UpdateUser(userID uuid.UUID, username, email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to get user from database: %v", err))
	}

	newUsername, newEmail, err := domain.DetermineUserDetailsToUpdate(user, username, email)
	if err != nil {
		return commonerrors.NewInvalidInput(fmt.Sprintf("invalid input: %v", err))
	}

	err = domain.ValidateUsername(newUsername)
	if err != nil {
		return commonerrors.NewInvalidInput(fmt.Sprintf("invalid username: %v", err))
	}

	err = domain.ValidateEmail(newEmail)
	if err != nil {
		return commonerrors.NewInvalidInput(fmt.Sprintf("invalid email: %v", err))
	}

	err = s.repo.UpdateUserDetails(ctx, userID, newUsername, newEmail)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to update user details: %v", err))
	}

	return nil
}

func (s *UserService) DeleteUser(userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.repo.DeleteUser(ctx, userID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to delete user: %v", err))
	}

	return nil
}
