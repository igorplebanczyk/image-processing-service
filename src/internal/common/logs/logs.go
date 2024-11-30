package logs

import (
	"log/slog"
	"os"
)

func init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("Init step 1 complete: logs initialized")
}
