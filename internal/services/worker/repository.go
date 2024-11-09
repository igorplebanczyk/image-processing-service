package worker

import "context"

type RefreshTokenRepository interface {
	DeleteExpiredRefreshTokens(ctx context.Context) error
}
