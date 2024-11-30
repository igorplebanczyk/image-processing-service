package tx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

type Provider struct {
	db *sql.DB
}

func NewProvider(db *sql.DB) *Provider {
	return &Provider{db: db}
}

func (p *Provider) Transact(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	queryErr := fn(tx)
	if queryErr == nil {
		return tx.Commit()
	}

	slog.Error("Database error: error performing transaction", "error", queryErr)

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		slog.Error("Database error: error rolling back transaction", "error", rollbackErr)
		return errors.Join(
			fmt.Errorf("error rolling back transaction: %w", rollbackErr),
			fmt.Errorf("error executing query: %w", queryErr),
		)
	}

	slog.Info("Transaction rolled back successfully")

	return queryErr
}
