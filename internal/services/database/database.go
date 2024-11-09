package database

import (
	"database/sql"
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
