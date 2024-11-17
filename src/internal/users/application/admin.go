package application

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"image-processing-service/src/internal/users/domain"
	"time"
)

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
