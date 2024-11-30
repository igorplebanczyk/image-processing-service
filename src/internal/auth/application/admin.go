package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	commonerrors "image-processing-service/src/internal/common/errors"
	"time"
)

func (s *AuthService) AdminLogoutUser(userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.refreshTokenDBRepo.RevokeAllUserRefreshTokens(ctx, userID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to revoke refresh tokens: %v", err))
	}

	return nil
}
