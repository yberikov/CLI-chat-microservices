package main

import (
	"chat/internal/app"
	"chat/internal/config"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	application := app.New(logger, cfg)

	logger.Info("starting server", slog.String("address", cfg.Address))

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := application.Run(ctx); err != nil {
			logger.Error("failed to start server", err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	logger.Info("stopping server: releasing all resources")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer shutdownCancel()
	if err := application.Stop(shutdownCtx); err != nil {
		logger.Error("failed to gracefully stop server", slog.String("error", err.Error()))
	}

	logger.Info("server gracefully stopped")
}
