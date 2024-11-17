package application

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"image-processing-service/src/internal/auth/domain"
	"time"
)

func (s *AuthService) AdminLogoutUser(userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.refreshTokenRepo.RevokeAllUserRefreshTokens(ctx, userID)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	return nil
}
