package server

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"hw3/internal/config"
	"hw3/internal/server/handlers"
	service2 "hw3/internal/services"
)

func New(cfg *config.Config, log *slog.Logger, service *service2.MessagerService) *http.Server {
	router := chi.NewRouter()

	router.HandleFunc("/getMessages", handlers.New(log, service))
	return &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	}
}
