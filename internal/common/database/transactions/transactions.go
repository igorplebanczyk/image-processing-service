package transactions

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type TransactionProvider struct {
	db *sql.DB
}

func NewTransactionProvider(db *sql.DB) *TransactionProvider {
	return &TransactionProvider{db: db}
}

func (p *TransactionProvider) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	queryErr := fn(tx)
	if queryErr == nil {
		return tx.Commit()
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		return errors.Join(
			fmt.Errorf("error rolling back transaction: %w", rollbackErr),
			fmt.Errorf("error executing query: %w", queryErr),
		)
	}

	return queryErr
}
