package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
)

func main() {
	slog.Info("Starting application...")

	app := &application{}
	err := app.assemble()
	if err != nil {
		slog.Error("Error configuring services", "error", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go app.dbWorker.Start()
	go app.serverService.Start()

	<-ctx.Done()
	slog.Info("Shutting down...")
	app.dbWorker.Stop()
	app.dbService.Stop()
	app.serverService.Stop()
}
