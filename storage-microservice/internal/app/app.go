package app

import (
	"context"
	"github.com/IBM/sarama"
	"hw3/internal/config"
	"hw3/internal/kafka"
	"hw3/internal/server"
	service2 "hw3/internal/services"
	"hw3/internal/storage/postgres"
	"log/slog"
	"net/http"
	"sync"
)

type App struct {
	log     *slog.Logger
	cfg     *config.Config
	server  *http.Server
	service *service2.MessagerService
}

func New(log *slog.Logger, cfg *config.Config) *App {
	storage, err := postgres.New(cfg.StoragePath)
	if err != nil {
		panic(err)
	}
	service := service2.NewService(storage)

	return &App{
		log:     log,
		cfg:     cfg,
		server:  server.New(cfg, log, service),
		service: service,
	}
}

func (a *App) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	consumer, err := kafka.RunConsumer(ctx, wg, a.cfg.Brokers, a.service)
	if err != nil {
		panic(err)
	}
	defer func(consumer sarama.ConsumerGroup) {
		err := consumer.Close()
		if err != nil {
			a.log.Error("error on closing consumer:", err)
		}
	}(consumer)

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
