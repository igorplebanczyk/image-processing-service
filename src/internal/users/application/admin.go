package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	commonerrors "image-processing-service/src/internal/common/errors"
	"image-processing-service/src/internal/users/domain"
	"time"
)

func (s *UserService) AdminGetAllUsers() ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, commonerrors.NewInternal(fmt.Sprintf("error fetching users from database: %v", err))
	}

	return users, nil
}

func (s *UserService) AdminUpdateUserRole(userID uuid.UUID, role domain.Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.repo.UpdateUserRole(ctx, userID, role)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error updating user role in database: %v", err))
	}

	return nil
}
