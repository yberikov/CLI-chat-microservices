package app

import (
	"chat/internal/config"
	"chat/internal/domain/models"
	"chat/internal/kafka"
	"chat/internal/server/handlers"
	"chat/internal/server/hub"
	service2 "chat/internal/service"
	"context"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"sync"
)

type App struct {
	log            *slog.Logger
	Server         *http.Server
	messageChannel chan models.Message
	cfg            *config.Config
}

func New(
	log *slog.Logger,
	cfg *config.Config,

) *App {
	messageChannel := make(chan models.Message)
	service := service2.NewService(messageChannel, cfg.StoragePath)
	router := chi.NewRouter()
	h := hub.NewHub(service, log)
	go h.Run()

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleWebSocket(h, w, r)
	})

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	}
	return &App{
		log:            log,
		Server:         srv,
		cfg:            cfg,
		messageChannel: messageChannel,
	}
}

func (a *App) Run(ctx context.Context) error {

	producer := kafka.NewProducer(a.log, a.cfg, a.messageChannel)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	producer.RunProducing(ctx, wg)

	go func() {
		if err := a.Server.ListenAndServe(); err != nil {
			a.log.Error("failed to start server", err)
		}
	}()
	wg.Wait()
	return nil
}

func (a *App) Stop(ctx context.Context) error {

	err := a.Server.Shutdown(ctx)
	if err != nil {
		return err
	}
	close(a.messageChannel)
	return nil
}
