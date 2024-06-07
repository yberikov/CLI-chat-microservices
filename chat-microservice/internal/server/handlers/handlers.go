package handlers

import (
	"chat/internal/domain/models"
	"chat/internal/server/hub"
	"github.com/gorilla/websocket"
	"log"
	"log/slog"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleWebSocket(h *hub.Hub, w http.ResponseWriter, r *http.Request) {
	log.Println("New connection request")
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

	for _, message := range h.Cache {
		toSend := message.Author + ": " + string(message.Text)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(toSend)); err != nil {
			h.Log.Error("Error on sending message of requesting username", err)
		}
	}
	client := &hub.Client{Hub: h, Conn: conn, Send: make(chan models.Message), Username: clientName}

	client.Hub.Register <- client

	go client.ReadMessage()
	go client.WriteMessage()
}
