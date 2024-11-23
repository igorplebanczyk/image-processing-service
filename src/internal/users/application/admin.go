package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	commonerrors "image-processing-service/src/internal/common/errors"
	"image-processing-service/src/internal/users/domain"
	"log/slog"
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

func (s *UserService) AdminBroadcast(subject, body string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error fetching users from database: %v", err))
	}

	errorChan := make(chan error, len(users))

	for _, user := range users {
		err = s.mailService.SendText([]string{user.Email}, subject, body)
		if err != nil {
			errorChan <- commonerrors.NewInternal(fmt.Sprintf("error sending email to user %s: %v", user.Email, err))
		}
	}

	close(errorChan)
	for err := range errorChan {
		if err != nil {
			slog.Error("Failed to send email", "error", err)
		}
	}

	return nil
}
