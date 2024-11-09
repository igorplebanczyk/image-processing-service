package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

type Service struct {
	DB *sql.DB
}

func New() *Service {
	return &Service{}
}

func (s *Service) Connect(user, password, host, port, dbName string) error {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}

	s.DB = db
	return nil
}

func (s *Service) withTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := s.DB.BeginTx(ctx, nil)
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
