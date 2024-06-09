package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"chat/internal/domain/models"
	"chat/internal/server/hub"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleWebSocket(h *hub.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.Log.Error("Error on upgrading connection", slog.String("err", err.Error()))
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, []byte("Enter your username: ")); err != nil {
		h.Log.Error("Error on sending message of requesting username", slog.String("err", err.Error()))
		return
	}

	_, username, err := conn.ReadMessage()
	if err != nil {
		h.Log.Error("Error on reading username:", slog.String("err", err.Error()))
		return
	}
	clientName := strings.TrimSpace(string(username))
	cache, err := h.Service.GetMessages()
	if err != nil {
		h.Log.Error("error on getting cache:", slog.String("err", err.Error()))
	}
	for _, message := range cache {
		toSend := message.Author + ": " + string(message.Text)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(toSend)); err != nil {
			h.Log.Error("Error on sending message of requesting username", slog.String("err", err.Error()))
		}
	}
	client := &hub.Client{Hub: h, Conn: conn, Send: make(chan models.Message), Username: clientName}

	client.Hub.Register <- client

	go client.ReadMessage()
	go client.WriteMessage()
}
