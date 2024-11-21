package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"image-processing-service/src/internal/common/logs"
	"log/slog"
)

type Service struct {
	db *sql.DB
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Connect(user, password, host, port, dbName string) (*sql.DB, error) {
	slog.Info("Connecting to database...", "type", logs.Standard)
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	s.db = db

	slog.Info("Connected to database")
	return db, nil
}

func (s *Service) Stop() {
	err := s.db.Close()
	if err != nil {
		slog.Error("Shutdown error: error closing database", "type", logs.Error, "error", err)
	}
	slog.Info("Database connection closed", "type", logs.Standard)
}
