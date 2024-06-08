package handlers

import (
	"encoding/json"
	service2 "hw3/internal/services"
	"log/slog"
	"net/http"
)

func New(log *slog.Logger, messager *service2.MessagerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		lastMessages, err := messager.GetLastMessages(10)
		if err != nil {
			http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
			log.Error("Failed to retrieve messages:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(lastMessages); err != nil {
			http.Error(w, "Failed to encode messages", http.StatusInternalServerError)
			log.Error("Failed to encode messages:", err)
			return
		}
	}
}
