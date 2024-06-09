package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hw3/internal/app"
	"hw3/internal/config"
)

func main() {
	cfg := config.MustLoad()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	application := app.New(logger, cfg)

	logger.Info("starting storage-server", slog.String("address", cfg.Address))

	ctx, cancel := context.WithCancel(context.Background())

	go application.Run(ctx)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	logger.Info("stopping server: releasing all resources")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, time.Second*5)
	defer shutdownCancel()
	if err := application.Stop(shutdownCtx); err != nil {
		logger.Error("failed to gracefully stop server", slog.String("error", err.Error()))
	}

	logger.Info("server gracefully stopped")
}
