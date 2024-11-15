package log

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func Setup(logDir string) error {
	currDate := time.Now().Format("2006-01-02")

	_, err := os.Stat(logDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0700)
		if err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}
	}

	logFileName := currDate + ".log"
	logFilePath := filepath.Clean(filepath.Join(logDir, logFileName))

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})
	logger := slog.New(handler)

	// Set the custom logger as the default logger globally
	slog.SetDefault(logger)

	return nil
}
