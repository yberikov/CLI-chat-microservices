package server

import (
	"net/http"

	"chat/internal/config"
	"chat/internal/server/handlers"
	"chat/internal/server/hub"
	"github.com/go-chi/chi/v5"
)

func New(cfg *config.Config, h *hub.Hub) *http.Server {
	router := chi.NewRouter()

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleWebSocket(h, w, r)
	})
	return &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	}
}
