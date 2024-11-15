package worker

import (
	"context"
	"database/sql"
	"image-processing-service/src/internal/common/database/transactions"
	"time"
)

type repository struct {
	db         *sql.DB
	txProvider *transactions.TransactionProvider
}

func newRepository(db *sql.DB, txProvider *transactions.TransactionProvider) repository {
	return repository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r repository) DeleteExpiredRefreshTokens(ctx context.Context) error {
	return r.txProvider.WithTransaction(ctx, func(tx *sql.Tx) error {
		now := time.Now()
		_, err := tx.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE expires_at < $1`, now)
		if err != nil {
			return err
		}

		return nil
	})
}
