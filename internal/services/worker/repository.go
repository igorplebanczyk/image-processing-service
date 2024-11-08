package worker

type RefreshTokenRepository interface {
	DeleteExpiredRefreshTokens() error
}
