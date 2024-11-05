package database

import (
	"database/sql"
	"fmt"
)

type Service struct {
	DB *sql.DB
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Connect(url string) error {
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
