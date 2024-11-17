package application

import (
	"context"
	"errors"
	"github.com/google/uuid"
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
		return nil, errors.Join(domain.ErrValidationFailed, err)
	}

	err = domain.ValidateEmail(email)
	if err != nil {
		return nil, errors.Join(domain.ErrValidationFailed, err)
	}

	err = domain.ValidatePassword(password)
	if err != nil {
		return nil, errors.Join(domain.ErrValidationFailed, err)
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.repo.CreateUser(ctx, username, email, hashedPassword)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	return user, nil
}

func (s *UserService) GetUser(userID uuid.UUID) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	return user, nil
}

func (s *UserService) UpdateUser(userID uuid.UUID, username, email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	newUsername, newEmail, err := domain.DetermineUserDetailsToUpdate(user, username, email)
	if err != nil {
		return errors.Join(domain.ErrInvalidRequest, err)
	}

	err = domain.ValidateUsername(newUsername)
	if err != nil {
		return errors.Join(domain.ErrValidationFailed, err)
	}

	err = domain.ValidateEmail(newEmail)
	if err != nil {
		return errors.Join(domain.ErrValidationFailed, err)
	}

	err = s.repo.UpdateUserDetails(ctx, userID, newUsername, newEmail)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	return nil
}

func (s *UserService) DeleteUser(userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.repo.DeleteUser(ctx, userID)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	return nil
}

func (s *UserService) AdminGetAllUsers() ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	return users, nil
}

func (s *UserService) AdminUpdateUserRole(userID uuid.UUID, role domain.Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.repo.UpdateUserRole(ctx, userID, role)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	return nil
}
