package logs

import (
	"log/slog"
	"os"
)

func Setup() error {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("Logger setup complete")

	return nil
}
