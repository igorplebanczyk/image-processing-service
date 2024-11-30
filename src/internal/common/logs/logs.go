package logs

import (
	"log/slog"
	"os"
)

// init() gets called implicitly via an import in the main package. slog.SetDefault() makes the custom logger the
// default globally, meaning calls like slog.Info() will use it automatically. Everything is printed out to the stdout
// so that it can be easily collected by Docker and then piped by promtail to Loki.

func init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("Init step 1 complete: logs initialized")
}
