package app

import (
	"chat/internal/config"
	"chat/internal/domain/models"
	"chat/internal/kafka"
	"chat/internal/server"
	"chat/internal/server/hub"
	service2 "chat/internal/service"
	"context"
	"log/slog"
	"net/http"
	"sync"
)

type App struct {
	log     *slog.Logger
	cfg     *config.Config
	kafkaCh chan models.Message
	hub     *hub.Hub
	server  *http.Server
}

func New(log *slog.Logger, cfg *config.Config) *App {
	kafkaCh := make(chan models.Message)
	service := service2.NewService(kafkaCh, cfg.StoragePath)
	h := hub.NewHub(log, service)

	return &App{
		log:     log,
		cfg:     cfg,
		kafkaCh: kafkaCh,
		hub:     h,
		server:  server.New(cfg, h),
	}
}

func (a *App) Run(ctx context.Context) {

	wg := &sync.WaitGroup{}

	wg.Add(1)
	producer := kafka.NewProducer(a.log, a.cfg, a.kafkaCh)
	go producer.RunProducing(ctx, wg)

	wg.Add(1)
	go a.hub.Run(ctx, wg)

	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			a.log.Info("serverRun: shutdown or closed")
		}
	}()
	wg.Wait()

}

func (a *App) Stop(ctx context.Context) error {
	ctx.Done()
	err := a.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
